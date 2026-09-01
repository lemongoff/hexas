package main

import (
	"context"
	"flag"
	"log"

	hexasconfig "github.com/lemongoff/hexas-config"
	"github.com/lemongoff/hexas/gateway"
)

var configFile = flag.String("f", "config/base.yaml", "the base YAML configuration file")

func main() {
	flag.Parse()

	manager, err := hexasconfig.NewManager(gateway.DefaultGatewayConf(),
		hexasconfig.YAMLFile(*configFile), hexasconfig.Environment("HEXAS_"))
	if err != nil {
		log.Fatal(err)
	}
	if err := manager.Load(context.Background()); err != nil {
		log.Fatal(err)
	}
	snapshot, ok := manager.Current()
	if !ok {
		log.Fatal("configuration was not published")
	}
	gw := gateway.MustNewServer(snapshot.Value())
	defer gw.Stop()
	gw.Start()
}
