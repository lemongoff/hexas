package proc

// DefaultShutdownConf returns Hexas process shutdown defaults.
func DefaultShutdownConf() ShutdownConf {
	return ShutdownConf{WrapUpTime: defaultWrapUpTime, WaitTime: defaultWaitTime}
}
