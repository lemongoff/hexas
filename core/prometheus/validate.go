package prometheus

import "fmt"

// Validate checks Prometheus endpoint configuration.
func (c Config) Validate() error {
	if c.Port < 0 || c.Port > 65535 {
		return fmt.Errorf("prometheus port must be in [0, 65535]")
	}
	if c.Port != 0 && c.Path == "" {
		return fmt.Errorf("prometheus path is required")
	}
	return nil
}
