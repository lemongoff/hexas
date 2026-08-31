package internal

import (
	"strings"

	"github.com/JellyGoFF/FF-Hexas/core/discov"
	"github.com/JellyGoFF/FF-Hexas/core/logx"
	"github.com/JellyGoFF/FF-Hexas/zrpc/resolver/internal/targets"
	"google.golang.org/grpc/resolver"
)

type discovBuilder struct{}

func (b *discovBuilder) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (
	resolver.Resolver, error) {
	hosts := strings.FieldsFunc(targets.GetHosts(target), func(r rune) bool {
		return r == EndpointSepChar
	})
	sub, err := discov.NewSubscriber(hosts, targets.GetKey(target))
	if err != nil {
		return nil, err
	}

	update := func() {
		vals := subset(sub.Values(), subsetSize)
		addrs := resolveAddresses(vals)
		if err := cc.UpdateState(resolver.State{
			Addresses: addrs,
		}); err != nil {
			logx.Error(err)
		}
	}
	sub.AddListener(update)
	update()

	return &discovResolver{
		cc:  cc,
		sub: sub,
	}, nil
}

func resolveAddresses(vals []string) []resolver.Address {
	addrs := make([]resolver.Address, 0, len(vals))
	for _, val := range vals {
		publishInfo, err := discov.DecodePublishInfo(val)
		if err != nil {
			logx.Errorf("DecodePublishInfo.Value: %s, err: %v", val, err)
			continue
		}
		if len(strings.TrimSpace(publishInfo.Addr)) == 0 {
			logx.Errorf("DecodePublishInfo.Value: %s, err: empty address", val)
			continue
		}

		logx.Infof("discovBuilder.Build serverName: %s, Addr: %s", publishInfo.ServerName, publishInfo.Addr)
		addrs = append(addrs, resolver.Address{
			Addr:       publishInfo.Addr,
			ServerName: publishInfo.ServerName,
		})
	}

	return addrs
}

func (b *discovBuilder) Scheme() string {
	return DiscovScheme
}

type discovResolver struct {
	cc  resolver.ClientConn
	sub *discov.Subscriber
}

func (r *discovResolver) Close() {
	r.sub.Close()
}

func (r *discovResolver) ResolveNow(_ resolver.ResolveNowOptions) {
}
