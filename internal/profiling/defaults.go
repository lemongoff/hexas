package profiling

// DefaultConfig returns Hexas continuous profiling defaults.
func DefaultConfig() Config {
	return Config{
		UploadRate: defaultUploadRate, CheckInterval: defaultCheckInterval,
		ProfilingDuration: defaultProfilingDuration, CpuThreshold: 700,
		ProfileType: ProfileType{CPU: true, Goroutines: true, Memory: true},
	}
}
