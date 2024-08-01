package org

import (
	"context"
	"encoding/json"
	"strconv"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GiverPowerLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGiverPowerLogic 分配算力
func NewGiverPowerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GiverPowerLogic {
	return &GiverPowerLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GiverPowerLogic) GiverPower(req *types.GivePowerReq) (resp *types.GivePowerResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.GivePowerResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	givePowerUserId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		l.Logger.Error("获取要暂停角色的用户ID失败", err)
		return &types.GivePowerResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 调用RPC接口 分配算力
	givePowerResp, err := l.svcCtx.OrgClient.GivePower(l.ctx, &user.GivePowerReq{
		UserId:          userId,
		GivePowerUserId: givePowerUserId,
		Power:           req.Power,
	})
	if err != nil {
		l.Logger.Error("调用RPC接口 分配算力失败", err)
		return &types.GivePowerResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	if givePowerResp.Code == consts.PowerNotEnough {
		return &types.GivePowerResp{
			BaseResp: types.BaseResp{
				Code: consts.PowerNotEnough,
				Msg:  "算力不足",
			},
		}, nil
	}

	return &types.GivePowerResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
	}, nil
}
