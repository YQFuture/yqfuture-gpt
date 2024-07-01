package shoptraininglogic

import (
	"context"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type TrainingGoodsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTrainingGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrainingGoodsLogic {
	return &TrainingGoodsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 训练商品
func (l *TrainingGoodsLogic) TrainingGoods(in *training.GoodsTrainingReq) (*training.GoodsTrainingResp, error) {
	// todo: add your logic here and delete this line

	return &training.GoodsTrainingResp{}, nil
}
