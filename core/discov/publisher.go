package discov

import (
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/lemongoff/hexas/core/discov/internal"
	"github.com/lemongoff/hexas/core/lang"
	"github.com/lemongoff/hexas/core/logc"
	"github.com/lemongoff/hexas/core/logx"
	"github.com/lemongoff/hexas/core/proc"
	"github.com/lemongoff/hexas/core/syncx"
	"github.com/lemongoff/hexas/core/threading"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type (
	// InstanceState describes whether a discovered instance accepts ordinary traffic.
	InstanceState string

	// PubOption defines the method to customize a Publisher.
	PubOption func(client *Publisher)

	PublishInfo struct {
		Addr       string
		InstanceID string
		State      InstanceState
	}

	// A Publisher can be used to publish the value to an etcd cluster on the given key.
	Publisher struct {
		endpoints  []string
		key        string
		fullKey    string
		id         int64
		value      string
		lease      clientv3.LeaseID
		quit       *syncx.DoneChan
		pauseChan  chan lang.PlaceholderType
		resumeChan chan lang.PlaceholderType
		mu         sync.Mutex
		client     internal.EtcdClient
	}
)

const (
	InstanceReady    InstanceState = "ready"
	InstanceDraining InstanceState = "draining"
)

func (s InstanceState) Valid() bool {
	return s == InstanceReady || s == InstanceDraining
}

// NewPublisher returns a Publisher.
// endpoints is the hosts of the etcd cluster.
// key:value are a pair to be published.
// opts are used to customize the Publisher.
func NewPublisher(endpoints []string, key, value string, opts ...PubOption) *Publisher {
	publisher := &Publisher{
		endpoints:  endpoints,
		key:        key,
		value:      value,
		quit:       syncx.NewDoneChan(),
		pauseChan:  make(chan lang.PlaceholderType),
		resumeChan: make(chan lang.PlaceholderType),
	}

	for _, opt := range opts {
		opt(publisher)
	}

	return publisher
}

// KeepAlive keeps key:value alive.
func (p *Publisher) KeepAlive() error {
	select {
	case <-p.quit.Done():
		return errors.New("etcd publisher is closed")
	default:
	}
	cli, err := p.doRegister()
	if err != nil {
		return err
	}

	proc.AddWrapUpListener(func() {
		p.Stop()
	})

	return p.keepAliveAsync(cli)
}

// Pause pauses the renewing of key:value.
func (p *Publisher) Pause() {
	p.pauseChan <- lang.Placeholder
}

// Resume resumes the renewing of key:value.
func (p *Publisher) Resume() {
	p.resumeChan <- lang.Placeholder
}

// Stop stops the renewing and revokes the registration.
func (p *Publisher) Stop() {
	_ = p.Close()
}

// Close stops renewing and synchronously revokes the active registration.
func (p *Publisher) Close() error {
	p.quit.Close()
	p.mu.Lock()
	client, lease := p.client, p.lease
	p.mu.Unlock()
	if client == nil || lease == clientv3.NoLease {
		return nil
	}
	return p.revoke(client, lease)
}

// Update atomically changes the value attached to the current lease.
func (p *Publisher) Update(value string) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.client == nil || p.lease == clientv3.NoLease || p.fullKey == "" {
		return errors.New("etcd publisher is not registered")
	}
	if _, err := p.client.Put(p.client.Ctx(), p.fullKey, value, clientv3.WithLease(p.lease)); err != nil {
		return err
	}
	p.value = value
	return nil
}

func (p *Publisher) doKeepAlive() error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-p.quit.Done():
			return nil
		default:
			cli, err := p.doRegister()
			if err != nil {
				logx.Errorf("etcd publisher doRegister: %v", err)
				break
			}

			if err := p.keepAliveAsync(cli); err != nil {
				logc.Errorf(cli.Ctx(), "etcd publisher keepAliveAsync: %v", err)
				break
			}

			return nil
		}
	}

	return nil
}

func (p *Publisher) doRegister() (internal.EtcdClient, error) {
	cli, err := internal.GetRegistry().GetConn(p.endpoints)
	if err != nil {
		return nil, err
	}

	_, err = p.register(cli)
	return cli, err
}

func (p *Publisher) keepAliveAsync(cli internal.EtcdClient) error {
	fullKey, value, lease := p.registration()
	ch, err := cli.KeepAlive(cli.Ctx(), lease)
	if err != nil {
		return err
	}

	threading.GoSafe(func() {
		wch := cli.Watch(cli.Ctx(), fullKey, clientv3.WithFilterPut())

		for {
			select {
			case _, ok := <-ch:
				if !ok {
					_ = p.revoke(cli, lease)
					if err := p.doKeepAlive(); err != nil {
						logc.Errorf(cli.Ctx(), "etcd publisher KeepAlive: %v", err)
					}
					return
				}

			case c, ok := <-wch:
				if !ok {
					_ = p.revoke(cli, lease)
					if err := p.doKeepAlive(); err != nil {
						logc.Errorf(cli.Ctx(), "etcd publisher KeepAlive: %v", err)
					}
					return
				}
				if c.Err() != nil {
					logc.Errorf(cli.Ctx(), "etcd publisher watch: %v", c.Err())
					_ = p.revoke(cli, lease)
					if err := p.doKeepAlive(); err != nil {
						logc.Errorf(cli.Ctx(), "etcd publisher KeepAlive: %v", err)
					}
					return
				}

				for _, evt := range c.Events {
					if evt.Type == clientv3.EventTypeDelete {
						logc.Infof(cli.Ctx(), "etcd publisher watch: %s, event: %v",
							evt.Kv.Key, evt.Type)
						_, currentValue, currentLease := p.registration()
						if currentLease != lease {
							continue
						}
						_, err := cli.Put(cli.Ctx(), fullKey, currentValue, clientv3.WithLease(lease))
						if err != nil {
							logc.Errorf(cli.Ctx(), "etcd publisher re-put key: %v", err)
						} else {
							logc.Infof(cli.Ctx(), "etcd publisher re-put key: %s, value: %s",
								fullKey, currentValue)
						}
					}
				}
			case <-p.pauseChan:
				logc.Infof(cli.Ctx(), "paused etcd renew, key: %s, value: %s", p.key, value)
				_ = p.revoke(cli, lease)
				select {
				case <-p.resumeChan:
					if err := p.doKeepAlive(); err != nil {
						logc.Errorf(cli.Ctx(), "etcd publisher KeepAlive: %v", err)
					}
					return
				case <-p.quit.Done():
					return
				}
			case <-p.quit.Done():
				_ = p.revoke(cli, lease)
				return
			}
		}
	})

	return nil
}

func (p *Publisher) register(client internal.EtcdClient) (clientv3.LeaseID, error) {
	resp, err := client.Grant(client.Ctx(), TimeToLive)
	if err != nil {
		return clientv3.NoLease, err
	}

	lease := resp.ID
	p.mu.Lock()
	value := p.value
	p.mu.Unlock()
	var fullKey string
	if p.id > 0 {
		fullKey = makeEtcdKey(p.key, p.id)
	} else {
		fullKey = makeEtcdKey(p.key, int64(lease))
	}
	_, err = client.Put(client.Ctx(), fullKey, value, clientv3.WithLease(lease))
	if err == nil {
		p.mu.Lock()
		p.client = client
		p.lease = lease
		p.fullKey = fullKey
		p.mu.Unlock()
	}

	return lease, err
}

func (p *Publisher) registration() (string, string, clientv3.LeaseID) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.fullKey, p.value, p.lease
}

func (p *Publisher) revoke(cli internal.EtcdClient, lease clientv3.LeaseID) error {
	if lease == clientv3.NoLease {
		return nil
	}
	p.mu.Lock()
	if p.lease != lease {
		p.mu.Unlock()
		return nil
	}
	p.lease = clientv3.NoLease
	p.client = nil
	p.mu.Unlock()
	_, err := cli.Revoke(cli.Ctx(), lease)
	if err != nil {
		logc.Errorf(cli.Ctx(), "etcd publisher revoke: %v", err)
	}
	return err
}

func EncodePublishInfo(info *PublishInfo) (string, error) {
	if info == nil {
		return "", errors.New("publish info is required")
	}
	if strings.TrimSpace(info.Addr) == "" {
		return "", errors.New("publish address is required")
	}
	if strings.TrimSpace(info.InstanceID) == "" {
		return "", errors.New("publish instance id is required")
	}
	if !info.State.Valid() {
		return "", errors.New("publish instance state is invalid")
	}
	buf, err := json.Marshal(info)
	if err != nil {
		return "", err
	}

	// data, err := base64.NewEncrypt().Encrypt(buf)
	// if err != nil {
	// 	return "", err
	// }

	return string(buf), nil
}

func DecodePublishInfo(data string) (*PublishInfo, error) {
	// buf, err := base64.NewEncrypt().Decrypt([]byte(data))
	// if err != nil {
	// 	return nil, err
	// }

	info := &PublishInfo{}
	if err := json.Unmarshal([]byte(data), info); err != nil {
		return nil, err
	}
	if strings.TrimSpace(info.Addr) == "" {
		return nil, errors.New("publish address is required")
	}
	if strings.TrimSpace(info.InstanceID) == "" {
		return nil, errors.New("publish instance id is required")
	}
	if !info.State.Valid() {
		return nil, errors.New("publish instance state is invalid")
	}

	return info, nil
}

// WithId customizes a Publisher with the id.
func WithId(id int64) PubOption {
	return func(publisher *Publisher) {
		publisher.id = id
	}
}

// WithPubEtcdAccount provides the etcd username/password.
func WithPubEtcdAccount(user, pass string) PubOption {
	return func(pub *Publisher) {
		RegisterAccount(pub.endpoints, user, pass)
	}
}

// WithPubEtcdTLS provides the etcd CertFile/CertKeyFile/CACertFile.
func WithPubEtcdTLS(certFile, certKeyFile, caFile string, insecureSkipVerify bool) PubOption {
	return func(pub *Publisher) {
		logx.Must(RegisterTLS(pub.endpoints, certFile, certKeyFile, caFile, insecureSkipVerify))
	}
}
