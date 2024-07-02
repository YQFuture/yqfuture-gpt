package mqs

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"time"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
)

type TrainingGoodsSuccess struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTrainingGoodsSuccess(ctx context.Context, svcCtx *svc.ServiceContext) *TrainingGoodsSuccess {
	return &TrainingGoodsSuccess{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TrainingGoodsSuccess) Consume(key, val string) error {
	logx.Infof("店铺训练商品消息消费成功 key :%s , val :%s", key, val)

	//TODO 消费消息，即发送商品消息给GPT进行训练

	//默认8个消费者，所以每次消费后延迟10秒，即每个消费者每分钟消费6条数据
	time.Sleep(time.Millisecond * time.Duration(l.svcCtx.Config.TrainingGoodsConf.ConsumeDelay))
	return nil
}
