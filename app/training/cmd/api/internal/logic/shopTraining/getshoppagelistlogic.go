package shopTraining

import (
	"context"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShopPageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetShopPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShopPageListLogic {
	return &GetShopPageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShopPageListLogic) GetShopPageList(req *types.ShopPageListReq) (resp *types.ShopPageListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
