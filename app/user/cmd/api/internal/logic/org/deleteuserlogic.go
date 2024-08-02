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

type DeleteUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewDeleteUserLogic 删除用户
func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteUserLogic) DeleteUser(req *types.DeleteUserReq) (resp *types.DeleteUserResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.DeleteUserResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	deleteUserId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		l.Logger.Error("获取要删除角色的用户ID失败", err)
		return &types.DeleteUserResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 调用RPC接口 删除用户
	_, err = l.svcCtx.OrgClient.DeleteUser(l.ctx, &org.DeleteUserReq{
		UserId:       userId,
		DeleteUserId: deleteUserId,
	})
	if err != nil {
		l.Logger.Error("调用RPC接口 删除用户失败", err)
		return &types.DeleteUserResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	return &types.DeleteUserResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
	}, nil
}
