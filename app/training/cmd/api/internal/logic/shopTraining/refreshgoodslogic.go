package shopTraining

import (
	"context"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshGoodsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewRefreshGoodsLogic 刷新商品
func NewRefreshGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshGoodsLogic {
	return &RefreshGoodsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefreshGoodsLogic) RefreshGoods(req *types.RefreshGoodsReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
