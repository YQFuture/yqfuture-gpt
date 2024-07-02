package mqs

import (
	"context"
	"yufuture-gpt/app/training/cmd/rpc/internal/config"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

func Consumers(c config.Config, ctx context.Context, svcContext *svc.ServiceContext) []service.Service {

	return []service.Service{
		//Listening for changes in consumption flow status
		kq.MustNewQueue(c.KqConsumerConf, NewTrainingGoodsSuccess(ctx, svcContext)),
		//.....
	}

}
