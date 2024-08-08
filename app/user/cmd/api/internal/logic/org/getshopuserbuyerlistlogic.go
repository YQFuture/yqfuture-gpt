package org

import (
	"context"
	"math/rand"
	"time"
	"yufuture-gpt/common/consts"

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

	return &types.ShopUserBuyerListResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
		Data: []types.ShopUserBuyer{
			{
				BuyerId:      "1",
				BuyerName:    "张三",
				BuyerHeadImg: "1bbd7a79-c5ec-4bb2-8453-65aa0f1631e2_OIP.jpg",
				StartTime:    time.Now().Unix(),
				AiReturnNum:  rand.Int63n(100),
				UsedPower:    rand.Int63n(1000000),
			},
		},
	}, nil
}
