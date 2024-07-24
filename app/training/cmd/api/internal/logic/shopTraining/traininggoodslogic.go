package shopTraining

import (
	"context"
	"encoding/json"
	"strconv"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/consts"

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

func (l *TrainingGoodsLogic) TrainingGoods(req *types.TrainingGoodsReq) (resp *types.TrainingGoodsResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return nil, err
	}
	goodsIdInt, err := strconv.ParseInt(req.GoodsId, 10, 64)
	if err != nil {
		l.Logger.Error("转换商品ID失败", err)
		return nil, err
	}
	_, err = l.svcCtx.ShopTrainingClient.TrainingGoods(l.ctx, &training.TrainingGoodsReq{
		UserId:  userId,
		GoodsId: goodsIdInt,
	})
	if err != nil {
		l.Logger.Error("训练商品失败", err)
		return nil, err
	}
	return &types.TrainingGoodsResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "开始训练商品成功",
		},
	}, nil
}
