package shoptraininglogic

import (
	"context"
	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"
	"yufuture-gpt/common/consts"

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
	tsShop, err := l.svcCtx.TsShopModel.FindOneByUuidAndUserId(l.ctx, in.UserId, in.Uuid)
	if err != nil {
		l.Logger.Error("根据uuid和userid查找店铺失败", err)
		return nil, err
	}
	shopTrainingProgressResp := &training.GetShopTrainingProgressResp{
		TrainingStatus: tsShop.TrainingStatus,
	}
	// 根据训练状态从不同数据源中找出数据
	// 预设完成 则从mongo中获取预设结果
	if tsShop.TrainingStatus == consts.PresettingComplete {
		shoppresettingshoptitles, err := l.svcCtx.ShoppresettingshoptitlesModel.FindNewOneByUuidAndUserId(l.ctx, in.Uuid, in.UserId)
		if err != nil {
			l.Logger.Error("从mongo中获取预设结果失败", err)
		}
		shopTrainingProgressResp.Token = shoppresettingshoptitles.PreSettingToken
		shopTrainingProgressResp.Power = shoppresettingshoptitles.PresettingPower
		shopTrainingProgressResp.FileSize = shoppresettingshoptitles.PresettingFileSize
	}
	// 训练完成 则从es中获取训练结果
	if tsShop.TrainingStatus == consts.TrainingComplete {

	}
	return shopTrainingProgressResp, nil
}
