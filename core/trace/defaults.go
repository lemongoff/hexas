package trace

// DefaultConfig returns Hexas tracing defaults.
func DefaultConfig() Config { return Config{Sampler: 1, Batcher: "otlpgrpc"} }
