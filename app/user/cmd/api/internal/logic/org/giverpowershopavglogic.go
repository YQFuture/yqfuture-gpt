package org

import (
	"context"
	"encoding/json"
	"yufuture-gpt/app/user/cmd/rpc/client/org"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GiverPowerShopAvgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGiverPowerShopAvgLogic 平均分配店铺算力
func NewGiverPowerShopAvgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GiverPowerShopAvgLogic {
	return &GiverPowerShopAvgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GiverPowerShopAvgLogic) GiverPowerShopAvg(req *types.GivePowerShopAvgReq) (resp *types.GivePowerShopAvgResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.GivePowerShopAvgResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 调用RPC接口 平均分配店铺算力
	_, err = l.svcCtx.OrgClient.GivePowerShopAvg(l.ctx, &org.GivePowerShopAvgReq{
		UserId: userId,
	})
	if err != nil {
		l.Logger.Error("调用RPC接口 平均分配店铺算力失败", err)
		return &types.GivePowerShopAvgResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	return &types.GivePowerShopAvgResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
	}, nil
}
