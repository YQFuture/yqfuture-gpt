package main

import (
	"flag"
	"fmt"

	"yufuture-gpt/app/training/cmd/api/internal/config"
	"yufuture-gpt/app/training/cmd/api/internal/handler"
	"yufuture-gpt/app/training/cmd/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

var configFile = flag.String("f", "C:\\Users\\ZXD\\Documents\\Code\\GoLang\\src\\yqfuture-gpt\\app\\training\\cmd\\api\\etc\\training.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
