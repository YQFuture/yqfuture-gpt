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

type GetOrgPermTreeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetOrgPermTreeLogic 获取团队权限树
func NewGetOrgPermTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgPermTreeLogic {
	return &GetOrgPermTreeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrgPermTreeLogic) GetOrgPermTree(req *types.OrgPermTreeReq) (resp *types.OrgPermTreeResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.OrgPermTreeResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	l.Logger.Info("获取用户ID成功", userId)

	// 调用RPC接口 获取团队权限列表
	orgPermListResp, err := l.svcCtx.OrgClient.GetOrgPermList(l.ctx, &user.OrgPermListReq{
		UserId: userId,
	})
	if err != nil {
		l.Logger.Error("获取团队权限列表失败", err)
		return &types.OrgPermTreeResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	l.Logger.Info("获取团队权限列表成功", orgPermListResp)

	// 团队权限列表解析成权限树
	var orgPermList []types.OrgPerm
	for _, perm := range orgPermListResp.Result {
		orgPermList = append(orgPermList, types.OrgPerm{
			PermId:   strconv.FormatInt(perm.PermId, 10),
			PermName: perm.PermName,
			PermCode: perm.PermCode,
			ParentId: strconv.FormatInt(perm.ParentId, 10),
		})
	}

	parentMap := make(map[string]*types.OrgPerm)
	for _, perm := range orgPermList {
		parentMap[perm.PermId] = &perm
	}
	var orgPermTree []types.OrgPerm
	for _, perm := range orgPermList {
		if perm.ParentId == "0" { // 假设根节点的ParentID为0
			orgPermTree = append(orgPermTree, perm)
		} else {
			parent, exists := parentMap[perm.ParentId]
			if exists {
				parent.Children = append(parent.Children, perm)
			}
		}
	}

	return &types.OrgPermTreeResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
		Data: orgPermTree,
	}, nil
}
