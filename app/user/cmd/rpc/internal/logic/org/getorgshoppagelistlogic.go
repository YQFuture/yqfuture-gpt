package orglogic

import (
	"context"
	"errors"
	model "yufuture-gpt/app/user/model/mongo"
	"yufuture-gpt/app/user/model/orm"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrgShopPageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrgShopPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgShopPageListLogic {
	return &GetOrgShopPageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetOrgShopPageList 获取团队店铺分页列表
func (l *GetOrgShopPageListLogic) GetOrgShopPageList(in *user.OrgShopPageListReq) (*user.OrgShopPageListResp, error) {
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
		return nil, errors.New("只有团队管理员才能获取团队店铺列表")
	}

	// 调用MongoDB获取团队权限文档
	dborgpermission, err := l.svcCtx.DborgpermissionModel.FindOne(l.ctx, bsOrg.MongoPermId)
	if err != nil {
		l.Logger.Error("获取团队权限文档失败: ", err)
		return nil, err
	}
	shopMap := make(map[int64]*model.ShopPerm)
	for _, shop := range dborgpermission.ShopPermList {
		shopMap[shop.Id] = shop
	}

	// 从MySQL中获取店铺列表
	orgShopListResult, err := l.svcCtx.BsShopModel.FindPageListByOrgId(l.ctx, bsOrg.Id, in.PageNum, in.PageSize, in.Query)
	if err != nil {
		l.Logger.Error("获取店铺列表失败: ", err)
		return nil, err
	}
	total, err := l.svcCtx.BsShopModel.FindPageTotalByOrgId(l.ctx, bsOrg.Id, in.PageNum, in.PageSize, in.Query)
	if err != nil {
		l.Logger.Error("获取店铺总数失败: ", err)
		return nil, err
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

	// 解析构建返回体
	orgShopPageListResp := &user.OrgShopPageListResp{
		PageNum:  in.PageNum,
		PageSize: in.PageSize,
		Total:    total,
	}
	var orgShopList []*user.OrgShop
	for _, shop := range *orgShopListResult {
		orgShop := &user.OrgShop{
			ShopId:          shop.Id,
			ShopName:        shop.ShopName,
			PlatformType:    shop.PlatformType,
			MonthPowerLimit: shop.MonthPowerLimit,
			MonthUsedPower:  shop.MonthUsedPower,
		}

		// 角色列表
		var shopRoleList []*user.ShopRole
		// 查找店铺关联的权限
		var shopPermId int64
		for _, perm := range dborgpermission.PermissionList {
			if perm.ResourceId == shop.Id && perm.Perm == "shop" {
				shopPermId = perm.Id
				break
			}
		}
		// 查找权限关联的角色
		for _, role := range dborgpermission.RoleList {
			for _, permId := range role.PermissionList {
				if *permId == shopPermId {
					shopRole := &user.ShopRole{
						RoleId:   role.Id,
						RoleName: role.Name,
					}
					shopRoleList = append(shopRoleList, shopRole)
					break
				}
			}
		}
		orgShop.RoleList = shopRoleList

		// 关键字转接用户列表
		var keywordSwitchingUserList []*user.ShopUser
		mongoShop := shopMap[shop.Id]
		if mongoShop != nil {
			for _, userId := range mongoShop.KeywordSwitchingUserList {
				keywordSwitchingUser := &user.ShopUser{
					UserId:   userMap[*userId].Id,
					NickName: userMap[*userId].NickName.String,
					Phone:    userMap[*userId].Phone.String,
					HeadImg:  userMap[*userId].HeadImg.String,
				}
				keywordSwitchingUserList = append(keywordSwitchingUserList, keywordSwitchingUser)
			}
			orgShop.KeywordSwitchingUserList = keywordSwitchingUserList

			// 异常责任人用户列表
			var exceptionDutyUserList []*user.ShopUser
			for _, userId := range mongoShop.ExceptionDutyUserList {
				exceptionResponsibleUser := &user.ShopUser{
					UserId:   userMap[*userId].Id,
					NickName: userMap[*userId].NickName.String,
					Phone:    userMap[*userId].Phone.String,
					HeadImg:  userMap[*userId].HeadImg.String,
				}
				exceptionDutyUserList = append(exceptionDutyUserList, exceptionResponsibleUser)
			}
			orgShop.ExceptionDutyUserList = exceptionDutyUserList
		}

		orgShopList = append(orgShopList, orgShop)
	}

	orgShopPageListResp.List = orgShopList
	return orgShopPageListResp, nil
}
