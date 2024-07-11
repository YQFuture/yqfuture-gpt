package shoptraininglogic

import (
	"context"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelPreSettingGoodsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelPreSettingGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelPreSettingGoodsLogic {
	return &CancelPreSettingGoodsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消预训练商品
func (l *CancelPreSettingGoodsLogic) CancelPreSettingGoods(in *training.CancelPreSettingGoodsReq) (*training.CancelPreSettingGoodsResp, error) {
	// todo: add your logic here and delete this line

	return &training.CancelPreSettingGoodsResp{}, nil
}
