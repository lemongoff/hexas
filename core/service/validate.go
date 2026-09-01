package service

import "fmt"

// Validate checks service bootstrap configuration.
func (c ServiceConf) Validate() error {
	if c.Mode != "" {
		switch c.Mode {
		case DevMode, TestMode, RtMode, PreMode, ProMode:
		default:
			return fmt.Errorf("invalid service mode %q", c.Mode)
		}
	}
	if err := c.Log.Validate(); err != nil {
		return fmt.Errorf("log: %w", err)
	}
	if err := c.Prometheus.Validate(); err != nil {
		return fmt.Errorf("prometheus: %w", err)
	}
	if err := c.Telemetry.Validate(); err != nil {
		return fmt.Errorf("telemetry: %w", err)
	}
	if err := c.DevServer.Validate(); err != nil {
		return fmt.Errorf("dev server: %w", err)
	}
	if c.Shutdown.WrapUpTime < 0 || c.Shutdown.WaitTime < 0 || (c.Shutdown.WrapUpTime > 0 && c.Shutdown.WaitTime > 0 && c.Shutdown.WaitTime <= c.Shutdown.WrapUpTime) {
		return fmt.Errorf("shutdown wait time must be greater than wrap-up time")
	}
	if c.Profiling.ServerAddr != "" {
		if err := c.Profiling.Validate(); err != nil {
			return fmt.Errorf("profiling: %w", err)
		}
	}
	return nil
}
