package cache

import "github.com/lemongoff/hexas/core/stores/redis"

// DefaultNodeConf returns Hexas cache node defaults.
func DefaultNodeConf() NodeConf {
	return NodeConf{RedisConf: redis.DefaultRedisConf(), Weight: 100}
}
