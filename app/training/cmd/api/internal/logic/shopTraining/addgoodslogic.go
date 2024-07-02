package shopTraining

import (
	"context"
	"encoding/json"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddGoodsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 添加商品
func NewAddGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddGoodsLogic {
	return &AddGoodsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddGoodsLogic) AddGoods(req *types.AddGoodsReq) (resp *types.BaseResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户id失败", err)
		return nil, err
	}
	_, err = l.svcCtx.ShopTrainingClient.AddGoods(l.ctx, &training.AddGoodsReq{
		ShopId: req.ShopId,
		UserId: userId,
		List:   req.List,
	})
	if err != nil {
		l.Logger.Error("添加商品失败", err)
		return nil, err
	}
	return &types.BaseResp{
		Code: consts.Success,
		Msg:  "添加商品成功",
	}, nil
}
