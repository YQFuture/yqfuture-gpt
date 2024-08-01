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

type GetOrgUserPageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetOrgUserPageListLogic 获取团队用户分页列表
func NewGetOrgUserPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgUserPageListLogic {
	return &GetOrgUserPageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrgUserPageListLogic) GetOrgUserPageList(req *types.OrgUserPageListReq) (resp *types.OrgUserPageListResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.OrgUserPageListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 调用RPC接口 获取团队用户分页列表
	orgUserPageResp, err := l.svcCtx.OrgClient.GetOrgUserPageList(l.ctx, &org.OrgUserPageListReq{
		UserId:   userId,
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Query:    req.Query,
	})
	if err != nil {
		l.Logger.Error("获取团队用户分页列表失败", err)
		return &types.OrgUserPageListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 封装返回数据
	resp = &types.OrgUserPageListResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
	}
	var orgUserPage types.OrgUserPage
	orgUserPage.BasePageResp.PageNum = orgUserPageResp.PageNum
	orgUserPage.BasePageResp.PageSize = orgUserPageResp.PageSize
	orgUserPage.BasePageResp.Total = orgUserPageResp.Total
	var orgUserList []types.OrgUser
	for _, orgUserResp := range orgUserPageResp.List {
		// 基础数据
		orgUser := types.OrgUser{
			UserId:   strconv.FormatInt(orgUserResp.UserId, 10),
			Phone:    orgUserResp.Phone,
			NickName: orgUserResp.NickName,
			HeadImg:  orgUserResp.HeadImg,
			Status:   orgUserResp.Status,
		}
		// 角色列表
		var roleList []types.UserRole
		for _, roleResp := range orgUserResp.RoleList {
			role := types.UserRole{
				RoleId:   strconv.FormatInt(roleResp.RoleId, 10),
				RoleName: roleResp.RoleName,
			}
			roleList = append(roleList, role)
		}
		orgUser.RoleList = roleList

		// 权限列表
		var permList []types.UserPerm
		for _, permResp := range orgUserResp.PermList {
			perm := types.UserPerm{
				PermId:   strconv.FormatInt(permResp.PermId, 10),
				PermName: permResp.PermName,
				PermCode: permResp.PermCode,
			}
			permList = append(permList, perm)
		}
		orgUser.PermList = permList

		orgUserList = append(orgUserList, orgUser)
	}

	orgUserPage.List = orgUserList
	resp.Data = orgUserPage
	return resp, nil
}
