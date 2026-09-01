package rest

import (
	"time"

	"github.com/lemongoff/hexas/core/service"
)

// DefaultMiddlewaresConf enables the standard REST middleware chain.
func DefaultMiddlewaresConf() MiddlewaresConf {
	return MiddlewaresConf{Trace: true, Log: true, Prometheus: true, MaxConns: true, Breaker: true, Shedding: true, Timeout: true, Recover: true, Metrics: true, MaxBytes: true, Gunzip: true}
}

// DefaultRestConf returns Hexas REST server defaults.
func DefaultRestConf() RestConf {
	return RestConf{
		ServiceConf: service.DefaultServiceConf(), Host: "0.0.0.0", MaxConns: 10000,
		MaxBytes: 1 << 20, Timeout: 3000, CpuThreshold: 900,
		Signature: SignatureConf{Expiry: time.Hour}, Middlewares: DefaultMiddlewaresConf(),
	}
}
