package orglogic

import (
	"context"
	"errors"
	"yufuture-gpt/app/user/model/orm"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetShopUserListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetShopUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShopUserListLogic {
	return &GetShopUserListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetShopUserList 获取店铺客服列表
func (l *GetShopUserListLogic) GetShopUserList(in *user.ShopUserListReq) (*user.ShopUserListResp, error) {
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
		return nil, errors.New("只有团队管理员才能获取获取店铺客服列表")
	}
	// 从MySQL中获取用户列表 并转换成map
	userListResult, err := l.svcCtx.BsUserModel.FindListByOrgId(l.ctx, bsOrg.Id)
	if err != nil {
		l.Logger.Error("获取用户列表失败: ", err)
		return nil, err
	}
	userMap := make(map[int64]*orm.BsUser)
	for _, userResult := range *userListResult {
		userMap[userResult.Id] = userResult
	}
	// 调用MongoDB获取团队权限文档
	dborgpermission, err := l.svcCtx.DborgpermissionModel.FindOne(l.ctx, bsOrg.MongoPermId)
	if err != nil {
		l.Logger.Error("获取团队权限文档失败: ", err)
		return nil, err
	}

	// 角色列表
	var shopRoleList []*user.ShopRole
	// 查找店铺关联的权限
	var shopPermId int64
	for _, perm := range dborgpermission.PermissionList {
		if perm.ResourceId == in.ShopId && perm.Perm == "shop" {
			shopPermId = perm.Id
			break
		}
	}
	// 查找权限关联的角色
	shopRoleMap := make(map[int64]*user.ShopRole)
	for _, role := range dborgpermission.RoleList {
		for _, permId := range role.PermissionList {
			if *permId == shopPermId {
				shopRole := &user.ShopRole{
					RoleId:   role.Id,
					RoleName: role.Name,
				}
				shopRoleList = append(shopRoleList, shopRole)
				shopRoleMap[shopRole.RoleId] = shopRole
				break
			}
		}
	}

	// 角色用户列表 即拥有该店铺权限的角色的用户列表
	var shopUserList []*user.ShopUser
	for _, mongoUser := range dborgpermission.UserList {
		for _, roleId := range mongoUser.RoleList {
			if shopRoleMap[*roleId] != nil {
				roleUser := &user.ShopUser{
					UserId:   userMap[mongoUser.Id].Id,
					NickName: userMap[mongoUser.Id].NickName.String,
					Phone:    userMap[mongoUser.Id].Phone.String,
					HeadImg:  userMap[mongoUser.Id].HeadImg.String,
				}
				careData, err := l.svcCtx.BsShopCareHistoryModel.FindOrgShopUserCareData(l.ctx, roleUser.UserId, in.ShopId, in.StartTime, in.EndTime)
				if err != nil {
					l.Logger.Error("获取店铺客服托管数据失败: ", err)
				} else {
					roleUser.CareTimes = careData.CareTimes
					roleUser.CareTime = careData.CareTime
					roleUser.UsedPower = careData.UsedPower
				}
				shopUserList = append(shopUserList, roleUser)
				break
			}
		}
	}

	return &user.ShopUserListResp{
		List: shopUserList,
	}, nil
}
