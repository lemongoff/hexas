package main

import (
	"flag"

	"github.com/JellyGoFF/FF-Hexas/core/conf"
	"github.com/JellyGoFF/FF-Hexas/gateway"
)

var configFile = flag.String("f", "etc/gateway.yaml", "config file")

func main() {
	flag.Parse()

	var c gateway.GatewayConf
	conf.MustLoad(*configFile, &c)
	gw := gateway.MustNewServer(c)
	defer gw.Stop()
	gw.Start()
}
