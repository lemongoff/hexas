package config

import "github.com/lemongoff/hexas/zrpc"

type Config struct{ zrpc.RpcServerConf }

func DefaultConfig() Config { return Config{RpcServerConf: zrpc.DefaultRpcServerConf()} }
