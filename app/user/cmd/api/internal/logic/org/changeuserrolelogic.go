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

type ChangeUserRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewChangeUserRoleLogic 修改用户角色
func NewChangeUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeUserRoleLogic {
	return &ChangeUserRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeUserRoleLogic) ChangeUserRole(req *types.ChangeUserRoleReq) (resp *types.ChangeUserRoleResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.ChangeUserRoleResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	changeUserId, err := strconv.ParseInt(req.UserId, 10, 64)
	if err != nil {
		l.Logger.Error("获取要修改角色的用户ID失败", err)
		return &types.ChangeUserRoleResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	var roleIdList []int64
	for _, v := range req.RoleIdList {
		roleId, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			l.Logger.Error("转换角色ID失败", err)
			return &types.ChangeUserRoleResp{
				BaseResp: types.BaseResp{
					Code: consts.Fail,
					Msg:  "操作失败",
				},
			}, nil
		}
		roleIdList = append(roleIdList, roleId)
	}

	// 调用RPC接口 修改用户角色
	_, err = l.svcCtx.OrgClient.ChangeUserRole(l.ctx, &user.ChangeUserRoleReq{
		UserId:       userId,
		ChangeUserId: changeUserId,
		RoleIdList:   roleIdList,
	})
	if err != nil {
		l.Logger.Error("调用RPC接口 修改用户角色失败", err)
		return &types.ChangeUserRoleResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	return &types.ChangeUserRoleResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
	}, nil
}
