package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"yufuture-gpt/app/training/cmd/api/internal/config"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
)

type ServiceContext struct {
	Config             config.Config
	ShopTrainingClient training.ShopTrainingClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		ShopTrainingClient: training.NewShopTrainingClient(zrpc.MustNewClient(c.TrainingClientConf).Conn()),
	}
}
