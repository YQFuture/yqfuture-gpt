package shoptraininglogic

import (
	"context"

	"yufuture-gpt/app/training/cmd/rpc/internal/svc"
	"yufuture-gpt/app/training/cmd/rpc/pb/training"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShopPageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetShopPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShopPageListLogic {
	return &GetShopPageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 查询店铺列表
func (l *GetShopPageListLogic) GetShopPageList(in *training.ShopPageListReq) (*training.ShopPageListResp, error) {
	// todo: add your logic here and delete this line

	return &training.ShopPageListResp{}, nil
}
