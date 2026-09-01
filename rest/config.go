package rest

import (
	"time"

	"github.com/lemongoff/hexas/core/service"
)

type (
	// MiddlewaresConf is the config of middlewares.
	MiddlewaresConf struct {
		Trace      bool
		Log        bool
		Prometheus bool
		MaxConns   bool
		Breaker    bool
		Shedding   bool
		Timeout    bool
		Recover    bool
		Metrics    bool
		MaxBytes   bool
		Gunzip     bool
	}

	// A PrivateKeyConf is a private key config.
	PrivateKeyConf struct {
		Fingerprint string
		KeyFile     string
	}

	// A SignatureConf is a signature config.
	SignatureConf struct {
		Strict      bool
		Expiry      time.Duration
		PrivateKeys []PrivateKeyConf
	}

	// A RestConf is a http service config.
	// Why not name it as Conf, because we need to consider usage like:
	//  type Config struct {
	//     zrpc.RpcConf
	//     rest.RestConf
	//  }
	// if with the name Conf, there will be two Conf inside Config.
	RestConf struct {
		service.ServiceConf
		Host     string
		Port     int
		CertFile string `json:",optional"`
		KeyFile  string `json:",optional"`
		Verbose  bool   `json:",optional"`
		MaxConns int
		MaxBytes int64
		// milliseconds
		Timeout      int64
		CpuThreshold int64
		Signature    SignatureConf `json:",optional"`
		// There are default values for all the items in Middlewares.
		Middlewares MiddlewaresConf
		// TraceIgnorePaths is paths blacklist for trace middleware.
		TraceIgnorePaths []string `json:",optional"`
	}
)
