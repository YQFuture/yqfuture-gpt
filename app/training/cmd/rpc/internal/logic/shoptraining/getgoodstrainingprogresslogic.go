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

// 获取商品训练进度
func (l *GetGoodsTrainingProgressLogic) GetGoodsTrainingProgress(in *training.GetGoodsTrainingProgressReq) (*training.GetGoodsTrainingProgressResp, error) {

	return &training.GetGoodsTrainingProgressResp{}, nil
}
