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

type TrainingShopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTrainingShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TrainingShopLogic {
	return &TrainingShopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TrainingShopLogic) TrainingShop(req *types.TrainingShopReq) (resp *types.BaseResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户id失败", err)
		return nil, err
	}
	_, err = l.svcCtx.ShopTrainingClient.TrainingShop(l.ctx, &training.TrainingShopReq{
		Uuid:   req.Uuid,
		UserId: userId,
	})
	if err != nil {
		l.Logger.Error("开启店铺训练失败", err)
		return nil, err
	}
	return &types.BaseResp{
		Code: consts.Success,
		Msg:  "success",
	}, nil
}
