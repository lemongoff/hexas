package trace

import "fmt"

// Validate checks tracing configuration invariants.
func (c Config) Validate() error {
	if c.Sampler < 0 || c.Sampler > 1 {
		return fmt.Errorf("trace sampler must be in [0, 1]")
	}
	if c.Batcher == "" {
		return nil
	}
	switch c.Batcher {
	case "zipkin", "otlpgrpc", "otlphttp", "file":
		return nil
	default:
		return fmt.Errorf("invalid trace batcher %q", c.Batcher)
	}
}
