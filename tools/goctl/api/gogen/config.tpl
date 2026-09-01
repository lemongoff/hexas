// Code scaffolded by goctl. Safe to edit.
// goctl {{.version}}

package config

import {{.authImport}}

type Config struct {
	rest.RestConf
	{{.auth}}
	{{.jwtTrans}}
}

func DefaultConfig() Config {
	return Config{RestConf: rest.DefaultRestConf()}
}
