package shoptraininglogic

import (
	"context"
	"yufuture-gpt/common/utils"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShopTrainingProgressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetShopTrainingProgressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShopTrainingProgressLogic {
	return &GetShopTrainingProgressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取店铺训练进度
func (l *GetShopTrainingProgressLogic) GetShopTrainingProgress(in *training.GetShopTrainingProgressReq) (*training.GetShopTrainingProgressResp, error) {
	//根据uuid和userid查找出店铺
	shop, err := l.svcCtx.TsShopModel.FindOneByUuidAndUserId(l.ctx, in.UserId, in.Uuid)
	if err != nil {
		l.Logger.Error("根据uuid和userid查找店铺失败", err)
		return nil, err
	}
	process, err := utils.AnyToString(shop)
	if err != nil {
		l.Logger.Info("序列化店铺训练进度失败", shop)
		return nil, err
	}
	return &training.GetShopTrainingProgressResp{
		Result: process,
	}, nil
}
