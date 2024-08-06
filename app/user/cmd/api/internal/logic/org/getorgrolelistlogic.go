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

type GetOrgRoleListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetOrgRoleListLogic 获取团队角色列表
func NewGetOrgRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgRoleListLogic {
	return &GetOrgRoleListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrgRoleListLogic) GetOrgRoleList(req *types.OrgRoleListReq) (resp *types.OrgRoleListResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.OrgRoleListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	//调用RPC接口 获取团队角色列表
	roleListResp, err := l.svcCtx.OrgClient.GetOrgRoleList(l.ctx, &org.OrgRoleListReq{
		UserId: userId,
		Query:  req.Query,
	})
	if err != nil {
		l.Logger.Error("获取团队角色列表失败: ", err)
		return &types.OrgRoleListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 解析返回的结果
	var orgRoleList []types.OrgRole
	for _, role := range roleListResp.Result {
		// 角色基本信息
		orgRole := types.OrgRole{
			RoleId:   strconv.FormatInt(role.RoleId, 10),
			RoleName: role.RoleName,
		}
		// 角色权限列表
		var rolePermList []types.RolePerm
		for _, perm := range role.PermList {
			rolePerm := types.RolePerm{
				PermId:   strconv.FormatInt(perm.PermId, 10),
				PermName: perm.PermName,
				PermCode: perm.PermCode,
			}
			rolePermList = append(rolePermList, rolePerm)
		}
		orgRole.PermList = rolePermList

		// 角色店铺列表
		var roleShopList []types.RoleShop
		for _, shop := range role.ShopList {
			roleShop := types.RoleShop{
				ShopId:       strconv.FormatInt(shop.ShopId, 10),
				ShopName:     shop.ShopName,
				PlatformType: shop.PlatformType,
			}
			roleShopList = append(roleShopList, roleShop)
		}
		orgRole.ShopList = roleShopList

		// 角色用户列表
		var roleUserList []types.RoleUser
		for _, user := range role.UserList {
			roleUser := types.RoleUser{
				UserId:   strconv.FormatInt(user.UserId, 10),
				NickName: user.NickName,
				HeadImg:  user.HeadImg,
			}
			roleUserList = append(roleUserList, roleUser)
		}
		orgRole.UserList = roleUserList

		orgRoleList = append(orgRoleList, orgRole)
	}

	return &types.OrgRoleListResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
		Data: orgRoleList,
	}, nil
}
