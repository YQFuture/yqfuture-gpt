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

type InviteJoinOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewInviteJoinOrgLogic 邀请用户加入团队
func NewInviteJoinOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteJoinOrgLogic {
	return &InviteJoinOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InviteJoinOrgLogic) InviteJoinOrg(req *types.InviteJoinOrgReq) (resp *types.InviteJoinOrgResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.InviteJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "邀请失败",
			},
		}, nil
	}
	inviteUserId, err := strconv.ParseInt(req.InviteUserId, 10, 64)
	if err != nil {
		l.Logger.Error("获取邀请用户ID失败", err)
		return &types.InviteJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "邀请失败",
			},
		}, nil
	}

	// 调用RPC接口 发起邀请用户加入团队请求
	inviteUserJoinOrgResp, err := l.svcCtx.OrgClient.InviteJoinOrg(l.ctx, &org.InviteJoinOrgReq{
		UserId:       userId,
		InviteUserId: inviteUserId,
	})
	if err != nil {
		l.Logger.Error("邀请用户加入团队失败", err)
		return &types.InviteJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "邀请失败",
			},
		}, nil
	}
	if inviteUserJoinOrgResp.Code == consts.OrgNumLimit {
		l.Logger.Error("用户加入团队达到上限", err)
		return &types.InviteJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.OrgNumLimit,
				Msg:  "用户加入团队达到上限",
			},
		}, nil
	}

	return &types.InviteJoinOrgResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "邀请成功",
		},
	}, nil
}
