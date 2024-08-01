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

type GiverPowerAvgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGiverPowerAvgLogic 平均分配算力
func NewGiverPowerAvgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GiverPowerAvgLogic {
	return &GiverPowerAvgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GiverPowerAvgLogic) GiverPowerAvg(req *types.GivePowerAvgReq) (resp *types.GivePowerAvgResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.GivePowerAvgResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 调用RPC接口 平均分配算力
	_, err = l.svcCtx.OrgClient.GivePowerAvg(l.ctx, &org.GivePowerAvgReq{
		UserId: userId,
	})
	if err != nil {
		l.Logger.Error("平均分配算力失败", err)
		return &types.GivePowerAvgResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	return &types.GivePowerAvgResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
	}, nil
}
