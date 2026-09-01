package serverinterceptors

// DefaultStatConf returns Hexas RPC stat interceptor defaults.
func DefaultStatConf() StatConf { return StatConf{SlowThreshold: defaultSlowThreshold} }
