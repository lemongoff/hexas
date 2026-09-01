package devserver

// Config is config for inner http server.
type Config struct {
	Enabled        bool
	Host           string `json:",optional"`
	Port           int
	MetricsPath    string
	HealthPath     string
	EnableMetrics  bool
	EnablePprof    bool
	HealthResponse string
}
