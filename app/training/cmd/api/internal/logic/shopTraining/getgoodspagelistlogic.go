package shopTraining

import (
	"context"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGoodsPageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetGoodsPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGoodsPageListLogic {
	return &GetGoodsPageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetGoodsPageListLogic) GetGoodsPageList(req *types.GoodsPageListReq) (resp *types.GoodsPageListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
