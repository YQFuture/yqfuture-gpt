package shoptraininglogic

import (
	"context"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnEnableGoodsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnEnableGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnEnableGoodsLogic {
	return &UnEnableGoodsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 禁用商品
func (l *UnEnableGoodsLogic) UnEnableGoods(in *training.GoodsTrainingReq) (*training.GoodsTrainingResp, error) {
	err := l.svcCtx.TsGoodsModel.UnEnableGoods(l.ctx, in)
	if err != nil {
		return nil, err
	}
	return &training.GoodsTrainingResp{}, nil
}
