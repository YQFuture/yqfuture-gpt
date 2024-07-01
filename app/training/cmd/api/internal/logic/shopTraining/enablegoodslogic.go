package shopTraining

import (
	"context"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type EnableGoodsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnableGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnableGoodsLogic {
	return &EnableGoodsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *EnableGoodsLogic) EnableGoods(req *types.BaseGoodsReq) (resp *types.BaseResp, err error) {
	_, err = l.svcCtx.ShopTrainingClient.EnableGoods(l.ctx, &training.GoodsTrainingReq{
		GoodsId: req.GoodsId,
	})
	if err != nil {
		l.Logger.Error("启用商品训练失败", err)
		return nil, err
	}
	return
}
