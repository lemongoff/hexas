package internal

import (
	"os"
	"strings"

	"github.com/lemongoff/hexas/core/discov"
	"github.com/lemongoff/hexas/core/netx"
)

const (
	allEths  = "0.0.0.0"
	envPodIp = "POD_IP"
)

// NewRpcPubServer returns a Server.
func NewRpcPubServer(etcd discov.EtcdConf, listenOn string,
	opts ...ServerOption) (Server, error) {
	pubListenOn := figureOutListenOn(listenOn)
	readyInfo, err := discov.EncodePublishInfo(&discov.PublishInfo{
		Addr:       pubListenOn,
		InstanceID: etcd.InstanceID,
		State:      discov.InstanceReady,
	})

	if err != nil {
		return nil, err
	}

	drainingInfo, err := discov.EncodePublishInfo(&discov.PublishInfo{
		Addr:       pubListenOn,
		InstanceID: etcd.InstanceID,
		State:      discov.InstanceDraining,
	})
	if err != nil {
		return nil, err
	}
	var pubOpts []discov.PubOption
	if etcd.HasAccount() {
		pubOpts = append(pubOpts, discov.WithPubEtcdAccount(etcd.User, etcd.Pass))
	}
	if etcd.HasTLS() {
		pubOpts = append(pubOpts, discov.WithPubEtcdTLS(etcd.CertFile, etcd.CertKeyFile,
			etcd.CACertFile, etcd.InsecureSkipVerify))
	}
	if etcd.HasID() {
		pubOpts = append(pubOpts, discov.WithId(etcd.ID))
	}
	publisher := discov.NewPublisher(etcd.Hosts, etcd.Key, readyInfo, pubOpts...)
	server := NewRpcServer(listenOn, opts...)
	server.SetLifecycle(
		publisher.KeepAlive,
		func() error { return publisher.Update(drainingInfo) },
		publisher.Close,
	)
	return server, nil
}

func figureOutListenOn(listenOn string) string {
	fields := strings.Split(listenOn, ":")
	if len(fields) == 0 {
		return listenOn
	}

	host := fields[0]
	if len(host) > 0 && host != allEths {
		return listenOn
	}

	ip := os.Getenv(envPodIp)
	if len(ip) == 0 {
		ip = netx.InternalIp()
	}
	if len(ip) == 0 {
		return listenOn
	}

	return strings.Join(append([]string{ip}, fields[1:]...), ":")
}
