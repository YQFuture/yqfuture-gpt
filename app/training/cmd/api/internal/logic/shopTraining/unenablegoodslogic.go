package shopTraining

import (
	"context"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnEnableGoodsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnEnableGoodsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnEnableGoodsLogic {
	return &UnEnableGoodsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnEnableGoodsLogic) UnEnableGoods(req *types.BaseGoodsReq) (resp *types.BaseResp, err error) {
	// todo: add your logic here and delete this line

	return
}
