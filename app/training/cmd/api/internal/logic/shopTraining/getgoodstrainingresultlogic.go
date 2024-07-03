package shopTraining

import (
	"context"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGoodsTrainingResultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取商品训练结果
func NewGetGoodsTrainingResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsTrainingResultLogic {
	return &GetGoodsTrainingResultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGoodsTrainingResultLogic) GetGoodsTrainingResult(req *types.GetGoodsTrainingResultReq) (resp *types.GetGoodsTrainingResultResp, err error) {
	// todo: add your logic here and delete this line

	return
}
