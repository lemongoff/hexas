package prometheus

// DefaultConfig returns Hexas Prometheus defaults.
func DefaultConfig() Config { return Config{Port: 9101, Path: "/metrics"} }
