package shopTraining

import (
	"context"
	"strconv"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/consts"

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

func (l *EnableGoodsLogic) EnableGoods(req *types.EnableGoodsReq) (resp *types.BaseResp, err error) {
	goodsIdInt, err := strconv.ParseInt(req.GoodsId, 10, 64)
	if err != nil {
		l.Logger.Error("转换商品id失败", err)
		return nil, err
	}
	_, err = l.svcCtx.ShopTrainingClient.EnableGoods(l.ctx, &training.EnableGoodsReq{
		GoodsId: goodsIdInt,
	})
	if err != nil {
		l.Logger.Error("启用商品训练失败", err)
		return nil, err
	}
	return &types.BaseResp{
		Code: consts.Success,
		Msg:  "启用商品训练成功",
	}, nil
}
