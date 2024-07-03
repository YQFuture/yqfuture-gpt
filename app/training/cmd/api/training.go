package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"yufuture-gpt/app/training/cmd/api/internal/types"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/training/cmd/api/internal/config"
	"yufuture-gpt/app/training/cmd/api/internal/handler"
	"yufuture-gpt/app/training/cmd/api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
)

var configFile = flag.String("f", "etc/training.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithUnauthorizedCallback(authFail))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

// authFail JWT认证失败自定义处理返回
func authFail(w http.ResponseWriter, r *http.Request, err error) {
	marshal, err := json.Marshal(&types.BaseResp{
		Code: consts.Unauthorized,
		Msg:  "权限不足",
	})
	if err != nil {
		return
	}
	_, err = w.Write(marshal)
	if err != nil {
		return
	}
}
