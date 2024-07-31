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

type AgreeInviteJoinOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewAgreeInviteJoinOrgLogic 用户同意邀请加入团队
func NewAgreeInviteJoinOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AgreeInviteJoinOrgLogic {
	return &AgreeInviteJoinOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AgreeInviteJoinOrgLogic) AgreeInviteJoinOrg(req *types.AgreeInviteJoinOrgReq) (resp *types.AgreeInviteJoinOrgResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.AgreeInviteJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	messageId, err := strconv.ParseInt(req.MessageId, 10, 64)
	if err != nil {
		l.Logger.Error("获取消息ID失败", err)
		return &types.AgreeInviteJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 调用RPC接口 同意邀请加入团队
	agreeInviteJoinOrgResp, err := l.svcCtx.OrgClient.AgreeInviteJoinOrg(l.ctx, &org.AgreeInviteJoinOrgReq{
		UserId:    userId,
		MessageId: messageId,
	})
	if err != nil {
		l.Logger.Error("同意邀请加入团队", err)
		return &types.AgreeInviteJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	if agreeInviteJoinOrgResp.Code == consts.OrgNumLimit {
		l.Logger.Error("用户加入团队达到上限", err)
		return &types.AgreeInviteJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.OrgNumLimit,
				Msg:  "用户加入团队达到上限",
			},
		}, nil
	}
	if agreeInviteJoinOrgResp.Code == consts.UserNumLimit {
		l.Logger.Error("团队加入的用户达到上限", err)
		return &types.AgreeInviteJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.OrgNumLimit,
				Msg:  "团队加入的用户达到上限",
			},
		}, nil
	}

	return &types.AgreeInviteJoinOrgResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
	}, nil
}
