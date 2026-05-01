package main

import (
	"flag"
	"fmt"

	"weclawbotnotify/services/ilink/internal/config"
	"weclawbotnotify/services/ilink/internal/server"
	"weclawbotnotify/services/ilink/internal/svc"
	"weclawbotnotify/services/ilink/pb/weclawbotnotify/services/ilink/ilink"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/ilink.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		ilink.RegisterIlinkServer(grpcServer, server.NewIlinkServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
