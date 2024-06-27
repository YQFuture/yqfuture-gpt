package shoptraininglogic

import (
	"context"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"

	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPreSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreSettingLogic {
	return &PreSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PreSettingLogic) PreSetting(in *training.ShopTrainingReq) (*training.ShopTrainingResp, error) {
	return &training.ShopTrainingResp{
		Result: "success",
	}, nil
}
