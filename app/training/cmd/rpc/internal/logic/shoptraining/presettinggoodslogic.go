package shoptraininglogic

import (
	"context"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type PreSettingGoodsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPreSettingGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PreSettingGoodsLogic {
	return &PreSettingGoodsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 预训练商品
func (l *PreSettingGoodsLogic) PreSettingGoods(in *training.PreSettingGoodsReq) (*training.PreSettingGoodsResp, error) {
	// todo: add your logic here and delete this line

	return &training.PreSettingGoodsResp{}, nil
}
