package org

import (
	"context"
	"encoding/json"
	"strconv"
	"yufuture-gpt/app/user/cmd/rpc/client/org"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GiverPowerShopLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGiverPowerShopLogic 分配店铺算力
func NewGiverPowerShopLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GiverPowerShopLogic {
	return &GiverPowerShopLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GiverPowerShopLogic) GiverPowerShop(req *types.GivePowerShopReq) (resp *types.GivePowerShopResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.GivePowerShopResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	shopId, err := strconv.ParseInt(req.ShopId, 10, 64)
	if err != nil {
		l.Logger.Error("获取店铺ID失败", err)
		return &types.GivePowerShopResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 调用RPC接口 分配店铺算力
	givePowerShopResp, err := l.svcCtx.OrgClient.GivePowerShop(l.ctx, &org.GivePowerShopReq{
		UserId: userId,
		ShopId: shopId,
		Power:  req.Power,
	})
	if err != nil {
		l.Logger.Error("调用RPC接口 分配店铺算力失败", err)
		return &types.GivePowerShopResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	if givePowerShopResp.Code == consts.PowerNotEnough {
		return &types.GivePowerShopResp{
			BaseResp: types.BaseResp{
				Code: consts.PowerNotEnough,
				Msg:  "算力不足",
			},
		}, nil
	}

	return &types.GivePowerShopResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
	}, nil
}
