package orglogic

import (
	"context"
	"errors"
	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateRoleLogic {
	return &UpdateRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateRole 更新角色
func (l *UpdateRoleLogic) UpdateRole(in *user.UpdateRoleReq) (*user.UpdateRoleResp, error) {
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
		return nil, errors.New("只有团队管理员才能修改角色")
	}
	// 调用MongoDB获取团队权限文档
	dborgpermission, err := l.svcCtx.DborgpermissionModel.FindOne(l.ctx, bsOrg.MongoPermId)
	if err != nil {
		l.Logger.Error("获取团队权限文档失败: ", err)
		return nil, err
	}

	// 获取要更新的角色
	roleList := dborgpermission.RoleList
	permissionList := make([]*int64, len(in.PermIdList))
	for i, v := range in.PermIdList {
		permissionList[i] = &v
	}
	for _, role := range roleList {
		if role.Id == in.RoleId {
			role.Name = in.RoleName
			role.PermissionList = permissionList
			break
		}
	}

	// 更新MongoDB中的团队权限文档
	_, err = l.svcCtx.DborgpermissionModel.Update(l.ctx, dborgpermission)
	if err != nil {
		l.Logger.Error("更新团队权限文档失败: ", err)
		return nil, err
	}

	return &user.UpdateRoleResp{}, nil
}
