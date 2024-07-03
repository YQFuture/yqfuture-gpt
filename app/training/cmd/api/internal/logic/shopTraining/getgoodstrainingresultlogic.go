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

type GetGoodsTrainingResultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取商品训练结果
func NewGetGoodsTrainingResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsTrainingResultLogic {
	return &GetGoodsTrainingResultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGoodsTrainingResultLogic) GetGoodsTrainingResult(req *types.GetGoodsTrainingResultReq) (resp *types.GetGoodsTrainingResultResp, err error) {
	goodsIdInt, err := strconv.ParseInt(req.GoodsId, 10, 64)
	if err != nil {
		l.Logger.Error("转换商品id失败", err)
		return nil, err
	}
	result, err := l.svcCtx.ShopTrainingClient.GetGoodsTrainingResult(l.ctx, &training.GetGoodsTrainingResultReq{
		GoodsId: goodsIdInt,
	})
	if err != nil {
		l.Logger.Error("获取商品训练结果失败", err)
		return nil, err
	}
	return &types.GetGoodsTrainingResultResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "获取商品训练结果成功",
		},
		Data: result.Result,
	}, nil
}
