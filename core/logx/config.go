package logx

type (
	// A LogConf is a logging config.
	LogConf struct {
		// ServiceName represents the service name.
		ServiceName string `json:",optional"`
		// Mode represents the logging mode, default is `console`.
		// console: log to console.
		// file: log to file.
		// volume: used in k8s, prepend the hostname to the log file name.
		Mode string
		// Encoding represents the encoding type, default is `json`.
		// json: json encoding.
		// plain: plain text encoding, typically used in development.
		Encoding string
		// TimeFormat represents the time format, default is `2006-01-02T15:04:05.000Z07:00`.
		TimeFormat string `json:",optional"`
		// Path represents the log file path, default is `logs`.
		Path string
		// Level represents the log level, default is `info`.
		Level string
		// MaxContentLength represents the max content bytes, default is no limit.
		MaxContentLength uint32 `json:",optional"`
		// Compress represents whether to compress the log file, default is `false`.
		Compress bool `json:",optional"`
		// Stat represents whether to log statistics, default is `true`.
		Stat bool
		// KeepDays represents how many days the log files will be kept. Default to keep all files.
		// Only take effect when Mode is `file` or `volume`, both work when Rotation is `daily` or `size`.
		KeepDays int `json:",optional"`
		// StackCooldownMillis represents the cooldown time for stack logging, default is 100ms.
		StackCooldownMillis int
		// MaxBackups represents how many backup log files will be kept. 0 means all files will be kept forever.
		// Only take effect when RotationRuleType is `size`.
		// Even though `MaxBackups` sets 0, log files will still be removed
		// if the `KeepDays` limitation is reached.
		MaxBackups int
		// MaxSize represents how much space the writing log file takes up. 0 means no limit. The unit is `MB`.
		// Only take effect when RotationRuleType is `size`
		MaxSize int
		// Rotation represents the type of log rotation rule. Default is `daily`.
		// daily: daily rotation.
		// size: size limited rotation.
		Rotation string
		// FileTimeFormat is the size-rotation filename time layout, defaulting to RFC3339.
		// On Windows the default replaces colons with hyphens. Custom layouts must be valid filenames.
		FileTimeFormat string `json:",optional"`
		// FieldKeys represents the field keys.
		FieldKeys fieldKeyConf `json:",optional"`
	}

	fieldKeyConf struct {
		// CallerKey represents the caller key.
		CallerKey string
		// ContentKey represents the content key.
		ContentKey string
		// DurationKey represents the duration key.
		DurationKey string
		// LevelKey represents the level key.
		LevelKey string
		// SpanKey represents the span key.
		SpanKey string
		// TimestampKey represents the timestamp key.
		TimestampKey string
		// TraceKey represents the trace key.
		TraceKey string
		// TruncatedKey represents the truncated key.
		TruncatedKey string
	}
)
