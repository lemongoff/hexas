package service

import (
	"github.com/lemongoff/hexas/core/logx"
	"github.com/lemongoff/hexas/core/proc"
	"github.com/lemongoff/hexas/core/prometheus"
	"github.com/lemongoff/hexas/core/trace"
	"github.com/lemongoff/hexas/internal/devserver"
	"github.com/lemongoff/hexas/internal/profiling"
)

// DefaultServiceConf returns Hexas service lifecycle defaults.
func DefaultServiceConf() ServiceConf {
	return ServiceConf{
		Log: logx.DefaultLogConf(), Mode: ProMode,
		Prometheus: prometheus.DefaultConfig(), Telemetry: trace.DefaultConfig(),
		DevServer: devserver.DefaultConfig(), Shutdown: proc.DefaultShutdownConf(),
		Profiling: profiling.DefaultConfig(),
	}
}
