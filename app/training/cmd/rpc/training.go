package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"

	"yufuture-gpt/app/training/cmd/rpc/internal/config"
	knowledgebasetrainingServer "yufuture-gpt/app/training/cmd/rpc/internal/server/knowledgebasetraining"
	shoptrainingServer "yufuture-gpt/app/training/cmd/rpc/internal/server/shoptraining"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "C:\\Users\\ZXD\\Documents\\Code\\GoLang\\src\\yqfuture-gpt\\app\\training\\cmd\\rpc\\etc\\training.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		training.RegisterShopTrainingServer(grpcServer, shoptrainingServer.NewShopTrainingServer(ctx))
		training.RegisterKnowledgeBaseTrainingServer(grpcServer, knowledgebasetrainingServer.NewKnowledgeBaseTrainingServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)

	// register service to consul
	_ = consul.RegisterService(c.ListenOn, c.Consul)

	s.Start()
}
