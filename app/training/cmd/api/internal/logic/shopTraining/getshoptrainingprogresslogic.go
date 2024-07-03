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

type GetShopTrainingProgressLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetShopTrainingProgressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShopTrainingProgressLogic {
	return &GetShopTrainingProgressLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShopTrainingProgressLogic) GetShopTrainingProgress(req *types.GetShopTrainingProgressReq) (resp *types.GetShopTrainingProgressResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户id失败", err)
		return nil, err
	}
	result, err := l.svcCtx.ShopTrainingClient.GetShopTrainingProgress(l.ctx, &training.GetShopTrainingProgressReq{
		UserId: userId,
		Uuid:   req.Uuid,
	})
	if err != nil {
		l.Logger.Error("获取店铺训练进度失败", err)
		return nil, err
	}
	return &types.GetShopTrainingProgressResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "获取店铺训练进度成功",
		},
		Data: result.Result,
	}, nil
}
