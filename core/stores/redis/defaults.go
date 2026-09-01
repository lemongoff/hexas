package redis

import "time"

// DefaultRedisConf returns Hexas Redis client defaults.
func DefaultRedisConf() RedisConf {
	return RedisConf{Type: NodeType, NonBlock: true, Protocol: 3, MaintNotifications: "disabled", PingTimeout: time.Second}
}
