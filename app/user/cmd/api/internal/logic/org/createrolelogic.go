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

type CreateRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewCreateRoleLogic 创建角色
func NewCreateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRoleLogic {
	return &CreateRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRoleLogic) CreateRole(req *types.CreateRoleReq) (resp *types.CreateRoleResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.CreateRoleResp{
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
			return &types.CreateRoleResp{
				BaseResp: types.BaseResp{
					Code: consts.Fail,
					Msg:  "操作失败",
				},
			}, nil
		}
		permIdList = append(permIdList, permId)
	}

	// 调用RPC接口 创建角色
	_, err = l.svcCtx.OrgClient.CreateRole(l.ctx, &user.CreateRoleReq{
		UserId:     userId,
		RoleName:   req.RoleName,
		PermIdList: permIdList,
	})
	if err != nil {
		l.Logger.Error("创建角色失败", err)
		return &types.CreateRoleResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	return &types.CreateRoleResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
	}, nil
}
