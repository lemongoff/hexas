// Code scaffolded by goctl. Safe to edit.
// goctl {{.version}}

package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	{{.importPackages}}
)

var configFile = flag.String("f", "config/base.yaml", "the base YAML configuration file")

func main() {
	flag.Parse()

	manager, err := hexasconfig.NewManager(config.DefaultConfig(),
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
	c := snapshot.Value()

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
