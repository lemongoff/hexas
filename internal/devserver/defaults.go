package devserver

// DefaultConfig returns Hexas development server defaults.
func DefaultConfig() Config {
	return Config{Enabled: true, Port: 6060, MetricsPath: "/metrics", HealthPath: "/healthz", EnableMetrics: true, EnablePprof: true, HealthResponse: "OK"}
}
