package orglogic

import (
	"context"
	"errors"
	model "yufuture-gpt/app/user/model/mongo"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateShopAssignLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateShopAssignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateShopAssignLogic {
	return &UpdateShopAssignLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateShopAssign 编辑店铺指派
func (l *UpdateShopAssignLogic) UpdateShopAssign(in *user.UpdateShopAssignReq) (*user.UpdateShopAssignResp, error) {
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
		return nil, errors.New("只有团队管理员才能编辑店铺指派")
	}
	// 调用MongoDB获取团队权限文档
	dborgpermission, err := l.svcCtx.DborgpermissionModel.FindOne(l.ctx, bsOrg.MongoPermId)
	if err != nil {
		l.Logger.Error("获取团队权限文档失败: ", err)
		return nil, err
	}

	var shopPerm *model.ShopPerm
	for _, shop := range dborgpermission.ShopPermList {
		if shop.Id == in.ShopId {
			shopPerm = shop
			break
		}
	}

	// 更新角色列表
	// 传进来的角色ID列表转map
	roleMap := make(map[int64]struct{})
	for _, roleId := range in.RoleList {
		roleMap[roleId] = struct{}{}
	}
	// 查找店铺关联的权限
	var shopPermId int64
	for _, perm := range dborgpermission.PermissionList {
		if perm.ResourceId == in.ShopId && perm.Perm == "shop" {
			shopPermId = perm.Id
			break
		}
	}
	if shopPermId == 0 {
		l.Logger.Error("店铺权限不存在")
		return nil, errors.New("店铺权限不存在")
	}
	// 遍历角色列表处理权限
	for _, role := range dborgpermission.RoleList {
		var permMap = make(map[int64]struct{})
		for _, permId := range role.PermissionList {
			permMap[*permId] = struct{}{}
		}
		// 判断当前角色关联的权限中是否包括当前店铺权限
		if _, ok := permMap[shopPermId]; ok {
			// 当前角色包含该店铺权限 但是传进来的角色列表不包括当前角色 清理当前店铺权限
			if _, ok := roleMap[role.Id]; !ok {
				for i, permId := range role.PermissionList {
					if *permId == shopPermId {
						role.PermissionList = append(role.PermissionList[:i], role.PermissionList[i+1:]...)
						break
					}
				}
			}
		} else {
			// 当前角色不包含该店铺权限 传进来的角色列表包括当前角色 添加当前店铺权限
			if _, ok := roleMap[role.Id]; ok {
				role.PermissionList = append(role.PermissionList, &shopPermId)
			}
		}

	}

	// 关键词切换用户列表
	var keywordSwitchingUserList []*int64
	for _, userId := range in.KeywordSwitchingUserList {
		keywordSwitchingUserList = append(keywordSwitchingUserList, &userId)
	}
	shopPerm.KeywordSwitchingUserList = keywordSwitchingUserList

	// 异常责任用户列表
	var exceptionDutyUserList []*int64
	for _, userId := range in.ExceptionDutyUserList {
		exceptionDutyUserList = append(exceptionDutyUserList, &userId)
	}
	shopPerm.ExceptionDutyUserList = exceptionDutyUserList

	// 更新MongoDB中的团队权限文档
	_, err = l.svcCtx.DborgpermissionModel.Update(l.ctx, dborgpermission)
	if err != nil {
		l.Logger.Error("更新团队权限文档失败: ", err)
		return nil, err
	}

	return &user.UpdateShopAssignResp{}, nil
}
