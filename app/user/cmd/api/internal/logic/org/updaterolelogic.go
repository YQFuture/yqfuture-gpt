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

type UpdateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateRoleLogic 更新角色
func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateRoleLogic) UpdateRole(req *types.UpdateRoleReq) (resp *types.UpdateRoleResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.UpdateRoleResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	roleId, err := strconv.ParseInt(req.RoleId, 10, 64)
	if err != nil {
		l.Logger.Error("转换角色ID失败", err)
		return &types.UpdateRoleResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	var permIdList []int64
	for _, v := range req.PermIdList {
		permId, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			l.Logger.Error("转换权限ID失败", err)
			return &types.UpdateRoleResp{
				BaseResp: types.BaseResp{
					Code: consts.Fail,
					Msg:  "操作失败",
				},
			}, nil
		}
		permIdList = append(permIdList, permId)
	}

	// 调用RPC接口 更新角色
	_, err = l.svcCtx.OrgClient.UpdateRole(l.ctx, &user.UpdateRoleReq{
		UserId:     userId,
		RoleId:     roleId,
		RoleName:   req.RoleName,
		PermIdList: permIdList,
	})
	if err != nil {
		l.Logger.Error("更新角色失败", err)
		return &types.UpdateRoleResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	return &types.UpdateRoleResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
	}, nil
}
