package shoptraininglogic

import (
	"context"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshGoodsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshGoodsLogic {
	return &RefreshGoodsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RefreshGoods 刷新商品
func (l *RefreshGoodsLogic) RefreshGoods(in *training.RefreshGoodsReq) (*training.RefreshGoodsResp, error) {
	// todo: add your logic here and delete this line

	return &training.RefreshGoodsResp{}, nil
}
