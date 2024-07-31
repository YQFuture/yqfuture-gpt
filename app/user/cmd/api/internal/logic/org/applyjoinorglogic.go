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

type ApplyJoinOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewApplyJoinOrgLogic 用户申请加入团队
func NewApplyJoinOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyJoinOrgLogic {
	return &ApplyJoinOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyJoinOrgLogic) ApplyJoinOrg(req *types.ApplyJoinOrgReq) (resp *types.ApplyJoinOrgResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.ApplyJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "申请失败",
			},
		}, nil
	}
	orgId, err := strconv.ParseInt(req.OrgId, 10, 64)
	if err != nil {
		l.Logger.Error("获取团队ID失败", err)
		return &types.ApplyJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "申请失败",
			},
		}, nil
	}

	// 调用RPC接口 发起用户申请加入团队请求
	applyJoinOrgResp, err := l.svcCtx.OrgClient.ApplyJoinOrg(l.ctx, &org.ApplyJoinOrgReq{
		UserId: userId,
		OrgId:  orgId,
	})
	if err != nil {
		l.Logger.Error("申请加入团队失败", err)
		return &types.ApplyJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "申请失败",
			},
		}, nil
	}
	if applyJoinOrgResp.Code == consts.OrgNumLimit {
		l.Logger.Error("用户加入团队达到上限", err)
		return &types.ApplyJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.OrgNumLimit,
				Msg:  "用户加入团队达到上限",
			},
		}, nil
	}
	if applyJoinOrgResp.Code == consts.UserNumLimit {
		l.Logger.Error("团队加入的用户达到上限", err)
		return &types.ApplyJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.OrgNumLimit,
				Msg:  "团队加入的用户达到上限",
			},
		}, nil
	}

	return &types.ApplyJoinOrgResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "申请成功",
		},
	}, nil
}
