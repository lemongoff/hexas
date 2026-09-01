package sqlx

// DefaultSqlConf returns Hexas SQL defaults.
func DefaultSqlConf() SqlConf { return SqlConf{DriverName: "mysql", Policy: "round-robin"} }
