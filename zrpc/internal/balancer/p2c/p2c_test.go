package p2c

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"

	"github.com/lemongoff/hexas/core/logx"
	"github.com/lemongoff/hexas/core/mathx"
	"github.com/lemongoff/hexas/core/stringx"
	"github.com/lemongoff/hexas/zrpc/route"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/status"
)

func init() {
	logx.Disable()
}

func TestP2cPicker_PickNil(t *testing.T) {
	builder := new(p2cPickerBuilder)
	picker := builder.Build(base.PickerBuildInfo{})
	_, err := picker.Pick(balancer.PickInfo{
		FullMethodName: "/",
		Ctx:            context.Background(),
	})
	assert.NotNil(t, err)
}

func TestP2cPicker_Pick(t *testing.T) {
	tests := []struct {
		name       string
		candidates int
		err        error
		threshold  float64
	}{
		{
			name:       "empty",
			candidates: 0,
			err:        balancer.ErrNoSubConnAvailable,
		},
		{
			name:       "single",
			candidates: 1,
			threshold:  0.9,
		},
		{
			name:       "two",
			candidates: 2,
			threshold:  0.5,
		},
		{
			name:       "multiple",
			candidates: 100,
			threshold:  0.95,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			const total = 10000
			builder := new(p2cPickerBuilder)
			ready := make(map[balancer.SubConn]base.SubConnInfo)
			for i := 0; i < test.candidates; i++ {
				ready[mockClientConn{
					id: stringx.Rand(),
				}] = base.SubConnInfo{
					Address: resolver.Address{
						Addr: strconv.Itoa(i),
					},
				}
			}

			picker := builder.Build(base.PickerBuildInfo{
				ReadySCs: ready,
			})
			var wg sync.WaitGroup
			wg.Add(total)
			for i := 0; i < total; i++ {
				result, err := picker.Pick(balancer.PickInfo{
					FullMethodName: "/",
					Ctx:            context.Background(),
				})
				assert.Equal(t, test.err, err)

				if test.err != nil {
					return
				}

				if i%100 == 0 {
					err = status.Error(codes.DeadlineExceeded, "deadline")
				}

				go func() {
					runtime.Gosched()
					result.Done(balancer.DoneInfo{
						Err: err,
					})
					wg.Done()
				}()
			}

			wg.Wait()
			dist := make(map[any]int)
			conns := picker.(*p2cPicker).conns
			for _, conn := range conns {
				dist[conn.addr.Addr] = int(conn.requests)
			}

			entropy := mathx.CalcEntropy(dist)
			assert.True(t, entropy > test.threshold, fmt.Sprintf("entropy is %f, less than %f",
				entropy, test.threshold))
		})
	}
}

func TestPickerWithEmptyConns(t *testing.T) {
	var picker p2cPicker
	_, err := picker.Pick(balancer.PickInfo{
		FullMethodName: "/",
		Ctx:            context.Background(),
	})
	assert.ErrorIs(t, err, balancer.ErrNoSubConnAvailable)
}

func TestP2cPickerRoutesToRequiredInstance(t *testing.T) {
	first := mockClientConn{id: "first"}
	second := mockClientConn{id: "second"}
	picker := new(p2cPickerBuilder).Build(base.PickerBuildInfo{ReadySCs: map[balancer.SubConn]base.SubConnInfo{
		first:  {Address: route.SetInstanceID(resolver.Address{Addr: "127.0.0.1:8001"}, "hall-1")},
		second: {Address: route.SetInstanceID(resolver.Address{Addr: "127.0.0.1:8002"}, "hall-2")},
	}})

	ctx := route.WithTarget(context.Background(), "hall-2", route.Require)
	result, err := picker.Pick(balancer.PickInfo{FullMethodName: "/", Ctx: ctx})
	assert.NoError(t, err)
	assert.Equal(t, second, result.SubConn)
	result.Done(balancer.DoneInfo{})
}

func TestP2cPickerExcludesDrainingFromOrdinaryTraffic(t *testing.T) {
	ready := mockClientConn{id: "ready"}
	draining := mockClientConn{id: "draining"}
	drainingAddress := route.SetInstanceID(resolver.Address{Addr: "127.0.0.1:8002"}, "hall-2")
	drainingAddress = route.SetDraining(drainingAddress, true)
	picker := new(p2cPickerBuilder).Build(base.PickerBuildInfo{ReadySCs: map[balancer.SubConn]base.SubConnInfo{
		ready:    {Address: route.SetInstanceID(resolver.Address{Addr: "127.0.0.1:8001"}, "hall-1")},
		draining: {Address: drainingAddress},
	}})

	result, err := picker.Pick(balancer.PickInfo{FullMethodName: "/", Ctx: context.Background()})
	assert.NoError(t, err)
	assert.Equal(t, ready, result.SubConn)
	result.Done(balancer.DoneInfo{})

	ctx := route.WithTarget(context.Background(), "hall-2", route.Require)
	result, err = picker.Pick(balancer.PickInfo{FullMethodName: "/", Ctx: ctx})
	assert.NoError(t, err)
	assert.Equal(t, draining, result.SubConn)
	result.Done(balancer.DoneInfo{})
}

func TestP2cPickerPreferDoesNotSelectDrainingInstance(t *testing.T) {
	ready := mockClientConn{id: "ready"}
	draining := mockClientConn{id: "draining"}
	drainingAddress := route.SetInstanceID(resolver.Address{Addr: "127.0.0.1:8002"}, "hall-2")
	drainingAddress = route.SetDraining(drainingAddress, true)
	picker := new(p2cPickerBuilder).Build(base.PickerBuildInfo{ReadySCs: map[balancer.SubConn]base.SubConnInfo{
		ready:    {Address: route.SetInstanceID(resolver.Address{Addr: "127.0.0.1:8001"}, "hall-1")},
		draining: {Address: drainingAddress},
	}})

	ctx := route.WithTarget(context.Background(), "hall-2", route.Prefer)
	result, err := picker.Pick(balancer.PickInfo{FullMethodName: "/", Ctx: ctx})
	assert.NoError(t, err)
	assert.Equal(t, ready, result.SubConn)
	result.Done(balancer.DoneInfo{})
}

func TestP2cPickerOnlyDrainingRequiresExplicitTarget(t *testing.T) {
	draining := mockClientConn{id: "draining"}
	address := route.SetInstanceID(resolver.Address{Addr: "127.0.0.1:8002"}, "hall-2")
	address = route.SetDraining(address, true)
	picker := new(p2cPickerBuilder).Build(base.PickerBuildInfo{ReadySCs: map[balancer.SubConn]base.SubConnInfo{
		draining: {Address: address},
	}})

	_, err := picker.Pick(balancer.PickInfo{FullMethodName: "/", Ctx: context.Background()})
	assert.ErrorIs(t, err, balancer.ErrNoSubConnAvailable)

	ctx := route.WithTarget(context.Background(), "hall-2", route.Require)
	result, err := picker.Pick(balancer.PickInfo{FullMethodName: "/", Ctx: ctx})
	assert.NoError(t, err)
	assert.Equal(t, draining, result.SubConn)
	result.Done(balancer.DoneInfo{})
}

func TestP2cPickerRequiredInstanceDoesNotFallback(t *testing.T) {
	conn := mockClientConn{id: "first"}
	picker := new(p2cPickerBuilder).Build(base.PickerBuildInfo{ReadySCs: map[balancer.SubConn]base.SubConnInfo{
		conn: {Address: route.SetInstanceID(resolver.Address{Addr: "127.0.0.1:8001"}, "hall-1")},
	}})

	ctx := route.WithTarget(context.Background(), "hall-missing", route.Require)
	_, err := picker.Pick(balancer.PickInfo{FullMethodName: "/", Ctx: ctx})
	assert.Equal(t, codes.Unavailable, status.Code(err))
}

func TestP2cPickerPreferInstanceFallsBack(t *testing.T) {
	conn := mockClientConn{id: "first"}
	picker := new(p2cPickerBuilder).Build(base.PickerBuildInfo{ReadySCs: map[balancer.SubConn]base.SubConnInfo{
		conn: {Address: route.SetInstanceID(resolver.Address{Addr: "127.0.0.1:8001"}, "hall-1")},
	}})

	ctx := route.WithTarget(context.Background(), "hall-missing", route.Prefer)
	result, err := picker.Pick(balancer.PickInfo{FullMethodName: "/", Ctx: ctx})
	assert.NoError(t, err)
	assert.Equal(t, conn, result.SubConn)
	result.Done(balancer.DoneInfo{})
}

func TestP2cPickerRejectsUnknownRouteMode(t *testing.T) {
	conn := mockClientConn{id: "first"}
	picker := new(p2cPickerBuilder).Build(base.PickerBuildInfo{ReadySCs: map[balancer.SubConn]base.SubConnInfo{
		conn: {Address: route.SetInstanceID(resolver.Address{Addr: "127.0.0.1:8001"}, "hall-1")},
	}})

	ctx := route.WithTarget(context.Background(), "hall-1", route.Mode(99))
	_, err := picker.Pick(balancer.PickInfo{FullMethodName: "/", Ctx: ctx})
	assert.Equal(t, codes.InvalidArgument, status.Code(err))
}

type mockClientConn struct {
	balancer.SubConn
	// add random string member to avoid map key equality.
	id string
}

func (m mockClientConn) GetOrBuildProducer(builder balancer.ProducerBuilder) (
	p balancer.Producer, close func()) {
	return builder.Build(m)
}

func (m mockClientConn) UpdateAddresses(_ []resolver.Address) {
}

func (m mockClientConn) Connect() {
}

func (m mockClientConn) Shutdown() {
}

func (m mockClientConn) RegisterHealthListener(func(balancer.SubConnState)) {
}
