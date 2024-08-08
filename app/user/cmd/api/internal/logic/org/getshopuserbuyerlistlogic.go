package org

import (
	"context"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShopUserBuyerListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetShopUserBuyerListLogic 获取店铺客服买家列表
func NewGetShopUserBuyerListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShopUserBuyerListLogic {
	return &GetShopUserBuyerListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShopUserBuyerListLogic) GetShopUserBuyerList(req *types.ShopUserBuyerListReq) (resp *types.ShopUserBuyerListResp, err error) {
	// todo: add your logic here and delete this line

	return
}
