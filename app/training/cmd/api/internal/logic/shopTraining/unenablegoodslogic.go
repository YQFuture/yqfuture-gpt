package shopTraining

import (
	"context"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnEnableGoodsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnEnableGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnEnableGoodsLogic {
	return &UnEnableGoodsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnEnableGoodsLogic) UnEnableGoods(req *types.BaseGoodsReq) (resp *types.BaseResp, err error) {
	_, err = l.svcCtx.ShopTrainingClient.UnEnableGoods(l.ctx, &training.GoodsTrainingReq{
		GoodsId: req.GoodsId,
	})
	if err != nil {
		l.Logger.Error("停用商品训练失败", err)
		return nil, err
	}
	return &types.BaseResp{
		Code: consts.Success,
		Msg:  "停用商品训练成功",
	}, nil
}
