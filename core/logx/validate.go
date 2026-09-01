package logx

import "fmt"

// Validate checks logging configuration invariants.
func (c LogConf) Validate() error {
	if c.Mode != "" && !oneOf(c.Mode, "console", "file", "volume") {
		return fmt.Errorf("invalid log mode %q", c.Mode)
	}
	if c.Encoding != "" && !oneOf(c.Encoding, "json", "plain") {
		return fmt.Errorf("invalid log encoding %q", c.Encoding)
	}
	if c.Level != "" && !oneOf(c.Level, "debug", "info", "error", "severe") {
		return fmt.Errorf("invalid log level %q", c.Level)
	}
	if c.Rotation != "" && !oneOf(c.Rotation, "daily", "size") {
		return fmt.Errorf("invalid log rotation %q", c.Rotation)
	}
	if c.StackCooldownMillis < 0 || c.MaxBackups < 0 || c.MaxSize < 0 || c.KeepDays < 0 {
		return fmt.Errorf("log limits must not be negative")
	}
	return nil
}

func oneOf(value string, allowed ...string) bool {
	for _, candidate := range allowed {
		if value == candidate {
			return true
		}
	}
	return false
}
