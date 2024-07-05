package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	// jwt相关配置
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
	TrainingClientConf zrpc.RpcClientConf
}
