package internal

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/lemongoff/hexas/core/discov"
	"github.com/lemongoff/hexas/zrpc/route"
	"github.com/stretchr/testify/assert"
	"go.etcd.io/etcd/client/v3/mock/mockserver"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

func TestDiscovBuilder_Scheme(t *testing.T) {
	var b discovBuilder
	assert.Equal(t, DiscovScheme, b.Scheme())
}

func TestResolveAddressesSkipsInvalidValues(t *testing.T) {
	first, err := discov.EncodePublishInfo(&discov.PublishInfo{Addr: "127.0.0.1:8001", InstanceID: "first"})
	assert.NoError(t, err)
	second, err := discov.EncodePublishInfo(&discov.PublishInfo{Addr: "127.0.0.1:8002", InstanceID: "second"})
	assert.NoError(t, err)

	addrs := resolveAddresses([]string{first, "not-json", `{}`, second})
	assert.Len(t, addrs, 2)
	assert.Equal(t, "127.0.0.1:8001", addrs[0].Addr)
	assert.Equal(t, "first", route.InstanceID(addrs[0]))
	assert.Equal(t, "127.0.0.1:8002", addrs[1].Addr)
	assert.Equal(t, "second", route.InstanceID(addrs[1]))
}

func TestDiscovBuilder_Build(t *testing.T) {
	servers, err := mockserver.StartMockServers(2)
	assert.NoError(t, err)
	t.Cleanup(func() {
		servers.Stop()
	})

	var addrs []string
	for _, server := range servers.Servers {
		addrs = append(addrs, server.Address)
	}
	u, err := url.Parse(fmt.Sprintf("%s:///%s?key=test", DiscovScheme, strings.Join(addrs, ",")))
	assert.NoError(t, err)

	var b discovBuilder
	_, err = b.Build(resolver.Target{
		URL: *u,
	}, mockClientConn{}, resolver.BuildOptions{})
	assert.Error(t, err)
}

type mockClientConn struct{}

func (m mockClientConn) UpdateState(_ resolver.State) error {
	return nil
}

func (m mockClientConn) ReportError(_ error) {
}

func (m mockClientConn) NewAddress(_ []resolver.Address) {
}

func (m mockClientConn) NewServiceConfig(_ string) {
}

func (m mockClientConn) ParseServiceConfig(_ string) *serviceconfig.ParseResult {
	return nil
}
