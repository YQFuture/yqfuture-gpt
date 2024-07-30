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

type AgreeApplyJoinOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewAgreeApplyJoinOrgLogic 同意用户申请加入团队
func NewAgreeApplyJoinOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AgreeApplyJoinOrgLogic {
	return &AgreeApplyJoinOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AgreeApplyJoinOrgLogic) AgreeApplyJoinOrg(req *types.AgreeApplyJoinOrgReq) (resp *types.AgreeApplyJoinOrgResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.AgreeApplyJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	messageId, err := strconv.ParseInt(req.MessageId, 10, 64)
	if err != nil {
		l.Logger.Error("获取消息ID失败", err)
		return &types.AgreeApplyJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 调用RPC接口 同意用户申请加入团队
	agreeApplyJoinOrgResp, err := l.svcCtx.OrgClient.AgreeApplyJoinOrg(l.ctx, &org.AgreeApplyJoinOrgReq{
		UserId:    userId,
		MessageId: messageId,
	})
	if err != nil {
		l.Logger.Error("同意用户申请加入团队失败", err)
		return &types.AgreeApplyJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	if agreeApplyJoinOrgResp.Code == consts.OrgNumLimit {
		l.Logger.Error("用户加入团队达到上限", err)
		return &types.AgreeApplyJoinOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.OrgNumLimit,
				Msg:  "用户加入团队达到上限",
			},
		}, nil
	}

	return &types.AgreeApplyJoinOrgResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
	}, nil
}
