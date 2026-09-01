package zrpc

import (
	"github.com/lemongoff/hexas/core/service"
	"github.com/lemongoff/hexas/core/stores/redis"
	"github.com/lemongoff/hexas/zrpc/internal"
)

// DefaultRpcClientConf returns Hexas RPC client defaults.
func DefaultRpcClientConf() RpcClientConf {
	return RpcClientConf{Timeout: 5000, BalancerName: "p2c_ewma", Middlewares: internal.DefaultClientMiddlewaresConf()}
}

// DefaultRpcServerConf returns Hexas RPC server defaults.
func DefaultRpcServerConf() RpcServerConf {
	return RpcServerConf{
		ServiceConf: service.DefaultServiceConf(), Redis: redis.RedisKeyConf{RedisConf: redis.DefaultRedisConf()},
		Timeout: 5000, CpuThreshold: 900, Health: true, Middlewares: internal.DefaultServerMiddlewaresConf(),
	}
}
