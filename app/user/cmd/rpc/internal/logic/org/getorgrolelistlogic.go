package orglogic

import (
	"context"
	"errors"
	"strings"
	"sync"
	model "yufuture-gpt/app/user/model/mongo"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrgRoleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrgRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgRoleListLogic {
	return &GetOrgRoleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetOrgRoleList 获取团队角色列表
func (l *GetOrgRoleListLogic) GetOrgRoleList(in *user.OrgRoleListReq) (*user.OrgRoleListResp, error) {
	// 获取当前用户数据和团队数据
	bsUser, err := l.svcCtx.BsUserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		l.Logger.Error("获取用户数据失败: ", err)
		return nil, err
	}
	bsOrg, err := l.svcCtx.BsOrganizationModel.FindOne(l.ctx, bsUser.NowOrgId)
	if err != nil {
		l.Logger.Error("获取团队数据失败: ", err)
		return nil, err
	}
	if bsOrg.OwnerId != bsUser.Id {
		l.Logger.Error("当前用户不是当前团队管理员")
		return nil, errors.New("只有团队管理员才能获取团队角色列表")
	}
	// 调用MongoDB获取团队权限文档
	dborgpermission, err := l.svcCtx.DborgpermissionModel.FindOne(l.ctx, bsOrg.MongoPermId)
	if err != nil {
		l.Logger.Error("获取团队权限文档失败: ", err)
		return nil, err
	}
	userList := dborgpermission.UserList
	roleList := dborgpermission.RoleList
	permList := dborgpermission.PermissionList
	permMap := make(map[int64]*model.Permission)
	for _, perm := range permList {
		permMap[perm.Id] = perm
	}

	var orgRoleList []*user.OrgRole

	// 获取角色列表后 并发获取每个角色相关的权限 店铺 用户数据
	var wg sync.WaitGroup
	wg.Add(len(roleList))
	for _, role := range roleList {
		go func(role *model.Role) {
			defer wg.Done()
			// 角色基本信息
			orgRole := &user.OrgRole{
				RoleId:   role.Id,
				RoleName: role.Name,
			}

			// 判断搜索条件
			if in.Query != "" {
				// 模糊搜索
				if !strings.Contains(role.Name, in.Query) {
					return
				}
			}

			// 角色权限列表
			var rolePermList []*user.RolePerm
			for _, permId := range role.PermissionList {
				perm := permMap[*permId]
				rolePerm := &user.RolePerm{
					PermId:   perm.Id,
					PermName: perm.Name,
					PermCode: perm.Perm,
				}
				rolePermList = append(rolePermList, rolePerm)
			}
			orgRole.PermList = rolePermList

			// 角色店铺列表

			// 角色用户列表
			var roleUserList []*user.RoleUser
			for _, roleUser := range userList {
				if roleUser.RoleList != nil && len(roleUser.RoleList) > 0 {
					for _, roleId := range roleUser.RoleList {
						if *roleId == role.Id {
							roleUser := &user.RoleUser{
								UserId: roleUser.Id,
							}
							roleUserList = append(roleUserList, roleUser)
						}
					}
				}
			}
			for _, roleUser := range roleUserList {
				roleBsUser, err := l.svcCtx.BsUserModel.FindOne(l.ctx, roleUser.UserId)
				if err != nil {
					l.Logger.Error("获取用户数据失败: ", err)
				}
				roleUser.NickName = roleBsUser.NickName.String
				roleUser.HeadImg = roleBsUser.HeadImg.String
			}
			orgRole.UserList = roleUserList

			orgRoleList = append(orgRoleList, orgRole)
		}(role)
	}
	wg.Wait()

	// 返回角色列表
	return &user.OrgRoleListResp{
		Result: orgRoleList,
	}, nil
}
