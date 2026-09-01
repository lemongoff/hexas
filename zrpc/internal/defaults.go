package internal

import "github.com/lemongoff/hexas/zrpc/internal/serverinterceptors"

// DefaultClientMiddlewaresConf enables the standard RPC client middleware chain.
func DefaultClientMiddlewaresConf() ClientMiddlewaresConf {
	return ClientMiddlewaresConf{Trace: true, Duration: true, Prometheus: true, Breaker: true, Timeout: true}
}

// DefaultServerMiddlewaresConf enables the standard RPC server middleware chain.
func DefaultServerMiddlewaresConf() ServerMiddlewaresConf {
	return ServerMiddlewaresConf{Trace: true, Recover: true, Stat: true, StatConf: serverinterceptors.DefaultStatConf(), Prometheus: true, Breaker: true}
}
