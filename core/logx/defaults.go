package logx

// DefaultLogConf returns Hexas logging defaults.
func DefaultLogConf() LogConf {
	return LogConf{
		Mode: "console", Encoding: "json", Path: "logs", Level: "info", Stat: true,
		StackCooldownMillis: 100, Rotation: "daily",
		FieldKeys: fieldKeyConf{
			CallerKey: "caller", ContentKey: "content", DurationKey: "duration",
			LevelKey: "level", SpanKey: "span", TimestampKey: "@timestamp",
			TraceKey: "trace", TruncatedKey: "truncated",
		},
	}
}
