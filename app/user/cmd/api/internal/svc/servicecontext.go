package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"yufuture-gpt/app/user/cmd/api/internal/config"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"
)

type ServiceContext struct {
	Config      config.Config
	Redis       *redis.Redis
	LoginClient user.LoginClient
	UserClient  user.UserClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := zrpc.MustNewClient(c.UserClientConf).Conn()
	return &ServiceContext{
		Config:      c,
		Redis:       redis.MustNewRedis(c.RedisConf),
		LoginClient: user.NewLoginClient(conn),
		UserClient:  user.NewUserClient(conn),
	}
}
