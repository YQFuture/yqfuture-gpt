package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"yufuture-gpt/app/user/cmd/api/internal/config"
	"yufuture-gpt/app/user/cmd/rpc/client/login"
	"yufuture-gpt/app/user/cmd/rpc/client/org"
	"yufuture-gpt/app/user/cmd/rpc/client/user"
)

type ServiceContext struct {
	Config      config.Config
	Redis       *redis.Redis
	LoginClient login.Login
	UserClient  user.User
	OrgClient   org.Org
}

func NewServiceContext(c config.Config) *ServiceContext {
	client := zrpc.MustNewClient(c.UserClientConf)
	return &ServiceContext{
		Config:      c,
		Redis:       redis.MustNewRedis(c.RedisConf),
		LoginClient: login.NewLogin(client),
		UserClient:  user.NewUser(client),
		OrgClient:   org.NewOrg(client),
	}
}
