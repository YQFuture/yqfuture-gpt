package org

import (
	"context"
	"encoding/json"
	"math/rand"
	"strconv"
	"time"
	"yufuture-gpt/app/user/cmd/rpc/client/org"
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
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.ShopUserBuyerListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	shopId, err := strconv.ParseInt(req.ShopId, 10, 64)
	if err != nil {
		l.Logger.Error("获取店铺ID失败", err)
		return &types.ShopUserBuyerListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 调用RPC接口 获取店铺客服买家列表
	l.Logger.Info(userId, shopId)
	userBuyerListResp, err := l.svcCtx.OrgClient.GetShopUserBuyerList(l.ctx, &org.ShopUserBuyerListReq{
		UserId:    userId,
		ShopId:    shopId,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Query:     req.Query,
	})
	if err != nil {
		l.Logger.Error("获取店铺客服买家列表失败", err)
		return &types.ShopUserBuyerListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	var data []types.ShopUserBuyer
	for _, userBuyer := range userBuyerListResp.List {
		data = append(data, types.ShopUserBuyer{
			BuyerId:      strconv.FormatInt(userBuyer.BuyerId, 10),
			BuyerName:    userBuyer.BuyerName,
			BuyerHeadImg: userBuyer.BuyerHeadImg,
			StartTime:    userBuyer.StartTime,
			AiReturnNum:  userBuyer.AiReturnNum,
			UsedPower:    userBuyer.UsedPower,
		})
	}
	/*	return &types.ShopUserBuyerListResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
		Data: data,
	}, nil*/

	return returnMockShopUserBuyerListResp()
}

func returnMockShopUserBuyerListResp() (resp *types.ShopUserBuyerListResp, err error) {
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
