package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	{{.imports}}

	hexasconfig "github.com/lemongoff/hexas-config"
	"github.com/lemongoff/hexas/core/service"
	"github.com/lemongoff/hexas/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
{{range .serviceNames}}       {{.Pkg}}.Register{{.GRPCService}}Server(grpcServer, {{.ServerPkg}}.New{{.Service}}Server(ctx))
{{end}}
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
