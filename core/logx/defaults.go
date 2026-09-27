package logx

import (
	"runtime"
	"time"
)

func defaultFileTimeFormat() string {
	if runtime.GOOS == "windows" {
		return "2006-01-02T15-04-05Z07-00"
	}
	return time.RFC3339
}

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
