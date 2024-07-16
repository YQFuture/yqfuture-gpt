package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"yufuture-gpt/app/training/cmd/rpc/internal/config"
	"yufuture-gpt/app/training/cmd/rpc/internal/mqs"
	basicfunctionServer "yufuture-gpt/app/training/cmd/rpc/internal/server/basicfunction"
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

var configFile = flag.String("f", "etc/training.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	server := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		training.RegisterShopTrainingServer(grpcServer, shoptrainingServer.NewShopTrainingServer(ctx))
		training.RegisterKnowledgeBaseTrainingServer(grpcServer, knowledgebasetrainingServer.NewKnowledgeBaseTrainingServer(ctx))
		training.RegisterBasicFunctionServer(grpcServer, basicfunctionServer.NewBasicFunctionServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)

	// 将服务注册到consul
	err := consul.RegisterService(c.ListenOn, c.Consul)
	if err != nil {
		panic(err)
	}

	// 创建服务组
	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()

	serviceGroup.Add(server)

	// 配置go-queue消息队列消费者
	for _, mq := range mqs.Consumers(c, context.Background(), ctx) {
		serviceGroup.Add(mq)
	}

	serviceGroup.Start()
}
