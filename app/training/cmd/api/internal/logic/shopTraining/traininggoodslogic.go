package shopTraining

import (
	"context"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TrainingGoodsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTrainingGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrainingGoodsLogic {
	return &TrainingGoodsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TrainingGoodsLogic) TrainingGoods(req *types.BaseGoodsReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
