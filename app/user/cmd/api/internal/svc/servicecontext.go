package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"yufuture-gpt/app/user/cmd/api/internal/config"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Redis:  redis.MustNewRedis(c.RedisConf),
	}
}
