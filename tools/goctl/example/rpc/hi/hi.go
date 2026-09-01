package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	hexasconfig "github.com/lemongoff/hexas-config"
	"github.com/lemongoff/hexas/core/service"
	"github.com/lemongoff/hexas/tools/goctl/example/rpc/hi/internal/config"
	eventServer "github.com/lemongoff/hexas/tools/goctl/example/rpc/hi/internal/server/event"
	greetServer "github.com/lemongoff/hexas/tools/goctl/example/rpc/hi/internal/server/greet"
	"github.com/lemongoff/hexas/tools/goctl/example/rpc/hi/internal/svc"
	"github.com/lemongoff/hexas/tools/goctl/example/rpc/hi/pb/hi"
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
		hi.RegisterGreetServer(grpcServer, greetServer.NewGreetServer(ctx))
		hi.RegisterEventServer(grpcServer, eventServer.NewEventServer(ctx))
		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
