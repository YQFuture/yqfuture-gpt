package shoptraininglogic

import (
	"context"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type EnableGoodsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEnableGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnableGoodsLogic {
	return &EnableGoodsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 启用商品
func (l *EnableGoodsLogic) EnableGoods(in *training.EnableGoodsReq) (*training.EnableGoodsResp, error) {
	err := l.svcCtx.TsGoodsModel.EnableGoods(l.ctx, in)
	if err != nil {
		l.Logger.Error("启用商品失败", err)
		return nil, err
	}
	return &training.EnableGoodsResp{}, nil
}
