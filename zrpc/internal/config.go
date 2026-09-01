package internal

import "github.com/lemongoff/hexas/zrpc/internal/serverinterceptors"

type (
	// StatConf defines the stat config.
	StatConf = serverinterceptors.StatConf

	// ClientMiddlewaresConf defines whether to use client middlewares.
	ClientMiddlewaresConf struct {
		Trace      bool
		Duration   bool
		Prometheus bool
		Breaker    bool
		Timeout    bool
	}

	// ServerMiddlewaresConf defines whether to use server middlewares.
	ServerMiddlewaresConf struct {
		Trace      bool
		Recover    bool
		Stat       bool
		StatConf   StatConf `json:",optional"`
		Prometheus bool
		Breaker    bool
	}

	// MethodTimeoutConf defines specified timeout for gRPC methods.
	MethodTimeoutConf = serverinterceptors.MethodTimeoutConf
)
