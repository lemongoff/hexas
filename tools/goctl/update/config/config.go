package config

import "github.com/lemongoff/hexas/core/logx"

// Config defines the goctl update service configuration.
type Config struct {
	logx.LogConf
	ListenOn string
	FileDir  string
	FilePath string
}

// DefaultConfig returns the update service defaults.
func DefaultConfig() Config { return Config{LogConf: logx.DefaultLogConf()} }
