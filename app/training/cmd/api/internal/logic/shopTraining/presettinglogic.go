package shopTraining

import (
	"context"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreSettingLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPreSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreSettingLogic {
	return &PreSettingLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PreSettingLogic) PreSetting(req *types.ShopTrainingReq) (resp *types.ShopTrainingResp, err error) {
	result, err := l.svcCtx.ShopTrainingClient.PreSetting(l.ctx, &training.ShopTrainingReq{})
	if err != nil {
		return nil, err
	}
	return &types.ShopTrainingResp{
		Data: result.Result,
		Code: 0,
		Msg:  "success",
	}, nil
}
