package proc

import "time"

const (
	defaultWrapUpTime = time.Second
	// Queues can block for five seconds, so allow them to finish before quitting.
	defaultWaitTime = 5500 * time.Millisecond
)

// ShutdownConf defines the process shutdown configuration on all platforms.
// Windows accepts this configuration, but does not implement timed shutdown.
type ShutdownConf struct {
	// WrapUpTime is the time to wait before calling shutdown listeners.
	WrapUpTime time.Duration
	// WaitTime is the time to wait before force quitting.
	WaitTime time.Duration
}

// DefaultShutdownConf returns Hexas process shutdown defaults.
func DefaultShutdownConf() ShutdownConf {
	return ShutdownConf{WrapUpTime: defaultWrapUpTime, WaitTime: defaultWaitTime}
}
