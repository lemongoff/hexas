package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	hexasconfig "github.com/lemongoff/hexas-config"
	"github.com/lemongoff/hexas/core/service"
	"github.com/lemongoff/hexas/tools/goctl/example/rpc/hello/internal/config"
	greetServer "github.com/lemongoff/hexas/tools/goctl/example/rpc/hello/internal/server/greet"
	"github.com/lemongoff/hexas/tools/goctl/example/rpc/hello/internal/svc"
	"github.com/lemongoff/hexas/tools/goctl/example/rpc/hello/pb/hello"
	"github.com/lemongoff/hexas/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "config/base.yaml", "the YAML configuration file")

func main() {
	flag.Parse()
	manager, err := hexasconfig.NewManager(config.DefaultConfig(), hexasconfig.YAMLFile(*configFile), hexasconfig.Environment("HEXAS_"))
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
		hello.RegisterGreetServer(grpcServer, greetServer.NewGreetServer(ctx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
