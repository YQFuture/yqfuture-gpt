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

type ResumeUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewResumeUserLogic 恢复用户
func NewResumeUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResumeUserLogic {
	return &ResumeUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResumeUserLogic) ResumeUser(req *types.ResumeUserReq) (resp *types.ResumeUserResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.ResumeUserResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	resumeUserId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		l.Logger.Error("获取要恢复角色的用户ID失败", err)
		return &types.ResumeUserResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 调用RPC接口 恢复用户
	_, err = l.svcCtx.OrgClient.ResumeUser(l.ctx, &org.ResumeUserReq{
		UserId:       userId,
		ResumeUserId: resumeUserId,
	})
	if err != nil {
		l.Logger.Error("调用RPC接口 恢复用户失败", err)
		return &types.ResumeUserResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	return &types.ResumeUserResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
	}, nil
}
