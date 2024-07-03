package shoptraininglogic

import (
	"context"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelPreSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelPreSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelPreSettingLogic {
	return &CancelPreSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消预训练
func (l *CancelPreSettingLogic) CancelPreSetting(in *training.CancelPreSettingReq) (*training.ShopTrainingResp, error) {
	shop, err := l.svcCtx.TsShopModel.FindOneByUuidAndUserId(l.ctx, in.UserId, in.Uuid)
	if err != nil {
		l.Logger.Error("根据uuid和userid查找店铺失败", err)
		return nil, err
	}
	l.Logger.Error("根据uuid和userid查找到店铺", shop)
	return &training.ShopTrainingResp{}, nil
}
