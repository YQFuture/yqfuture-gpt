package orglogic

import (
	"context"
	"errors"
	"yufuture-gpt/app/user/cmd/rpc/client/org"
	model "yufuture-gpt/app/user/model/mongo"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrgUserStatisticsPageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrgUserStatisticsPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgUserStatisticsPageListLogic {
	return &GetOrgUserStatisticsPageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetOrgUserStatisticsPageList 获取组织用户统计信息分页列表
func (l *GetOrgUserStatisticsPageListLogic) GetOrgUserStatisticsPageList(in *user.OrgUserStatisticsPageListReq) (*user.OrgUserStatisticsPageListResp, error) {
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
		return nil, errors.New("只有团队管理员才能获取团队用户列表")
	}

	// 调用MongoDB获取团队权限文档
	dborgpermission, err := l.svcCtx.DborgpermissionModel.FindOne(l.ctx, bsOrg.MongoPermId)
	if err != nil {
		l.Logger.Error("获取团队权限文档失败: ", err)
		return nil, err
	}
	permMap := make(map[int64]*model.Permission)
	for _, mongoPerm := range dborgpermission.PermissionList {
		permMap[mongoPerm.Id] = mongoPerm
	}
	// 从MySQL中获取用户列表
	userListResult, err := l.svcCtx.BsUserModel.FindPageListByOrgId(l.ctx, bsOrg.Id, in.PageNum, in.PageSize, in.Query)
	if err != nil {
		l.Logger.Error("获取团队用户列表失败: ", err)
		return nil, err
	}
	total, err := l.svcCtx.BsUserModel.FindPageTotalByOrgId(l.ctx, bsOrg.Id, in.PageNum, in.PageSize, in.Query)
	if err != nil {
		l.Logger.Error("获取团队用户总数失败: ", err)
		return nil, err
	}

	// 解析构建返回体
	var orgUserList []*user.OrgUserStatistics
	for _, userResult := range *userListResult {
		_, userPermList := GetUserRolePermList(userResult.Id, dborgpermission)
		var roleShopList []*user.RoleShop
		orgUser := &org.OrgUserStatistics{
			UserId:   userResult.Id,
			Phone:    userResult.Phone.String,
			NickName: userResult.NickName.String,
			HeadImg:  userResult.HeadImg.String,
		}

		for _, userPerm := range userPermList {
			perm := permMap[userPerm.PermId]
			// 用户店铺列表
			if perm.Perm == "shop" {
				bsShop, err := l.svcCtx.BsShopModel.FindOne(l.ctx, perm.ResourceId)
				if err != nil {
					l.Logger.Error("获取店铺数据失败: ", err)
					break
				}
				roleShop := &user.RoleShop{
					ShopId:       bsShop.Id,
					ShopName:     bsShop.ShopName,
					PlatformType: bsShop.PlatformType,
				}
				roleShopList = append(roleShopList, roleShop)
			}
		}
		orgUser.ShopList = roleShopList

		orgUserList = append(orgUserList, orgUser)
	}

	return &user.OrgUserStatisticsPageListResp{
		PageNum:  in.PageNum,
		PageSize: in.PageSize,
		Total:    total,
		List:     orgUserList,
	}, nil
}
