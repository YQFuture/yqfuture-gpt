package shopTraining

import (
	"context"

	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShopListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetShopListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShopListLogic {
	return &GetShopListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShopListLogic) GetShopList(req *types.ShopPageListReq) (resp *types.ShopPageListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
