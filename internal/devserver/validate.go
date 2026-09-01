package devserver

import "fmt"

// Validate checks development server configuration.
func (c Config) Validate() error {
	if c == (Config{}) {
		return nil
	}
	if c.Port < 0 || c.Port > 65535 {
		return fmt.Errorf("dev server port must be in [0, 65535]")
	}
	if c.MetricsPath == "" || c.HealthPath == "" {
		return fmt.Errorf("dev server paths are required")
	}
	return nil
}
