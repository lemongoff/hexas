package rest

import "fmt"

// Validate checks REST bootstrap configuration.
func (c RestConf) Validate() error {
	if err := c.ServiceConf.Validate(); err != nil {
		return err
	}
	if c.Port < 0 || c.Port > 65535 {
		return fmt.Errorf("rest port must be in [0, 65535]")
	}
	if c.MaxConns < 0 || c.MaxBytes < 0 || c.Timeout < 0 {
		return fmt.Errorf("rest limits and timeout must be positive")
	}
	if c.CpuThreshold < 0 || c.CpuThreshold >= 1000 {
		return fmt.Errorf("rest cpu threshold must be in [0, 1000)")
	}
	if c.Signature.Expiry < 0 {
		return fmt.Errorf("signature expiry must be positive")
	}
	return nil
}
