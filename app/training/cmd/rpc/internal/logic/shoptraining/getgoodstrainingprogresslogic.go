package shoptraininglogic

import (
	"context"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGoodsTrainingProgressLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGoodsTrainingProgressLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsTrainingProgressLogic {
	return &GetGoodsTrainingProgressLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGoodsTrainingProgress 获取商品训练进度
func (l *GetGoodsTrainingProgressLogic) GetGoodsTrainingProgress(in *training.GetGoodsTrainingProgressReq) (*training.GetGoodsTrainingProgressResp, error) {
	tsGoods, err := l.svcCtx.TsGoodsModel.FindOne(l.ctx, in.GoodsId)
	if err != nil {
		l.Logger.Error("根据商品ID查找商品失败", err)
		return nil, err
	}
	return &training.GetGoodsTrainingProgressResp{
		TrainingStatus: tsGoods.TrainingStatus,
	}, nil
}
