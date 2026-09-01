package profiling

import "fmt"

// Validate checks continuous profiling configuration.
func (c Config) Validate() error {
	if c.UploadRate <= 0 || c.CheckInterval <= 0 || c.ProfilingDuration <= 0 { return fmt.Errorf("profiling durations must be positive") }
	if c.CpuThreshold < 0 || c.CpuThreshold >= 1000 { return fmt.Errorf("profiling cpu threshold must be in [0, 1000)") }
	return nil
}
