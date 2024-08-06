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

type GetOrgUserPageListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrgUserPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgUserPageListLogic {
	return &GetOrgUserPageListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetOrgUserPageList 获取团队用户分页列表
func (l *GetOrgUserPageListLogic) GetOrgUserPageList(in *user.OrgUserPageListReq) (*user.OrgUserPageListResp, error) {
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
	var orgUserList []*org.OrgUser
	for _, userResult := range *userListResult {
		userRoleList, userPermList := GetUserRolePermList(userResult.Id, dborgpermission)
		orgUser := &org.OrgUser{
			UserId:          userResult.Id,
			Phone:           userResult.Phone.String,
			NickName:        userResult.NickName.String,
			HeadImg:         userResult.HeadImg.String,
			Status:          userResult.Status,
			MonthPowerLimit: userResult.MonthPowerLimit,
			MonthUsedPower:  userResult.MonthUsedPower,
			RoleList:        userRoleList,
			PermList:        userPermList,
		}
		orgUserList = append(orgUserList, orgUser)
	}

	// 获取当前团队已分配的总算力 判断剩余算力是否足够
	totalPower, err := l.svcCtx.BsUserOrgModel.FindOrgTotalGivePower(l.ctx, bsOrg.Id)
	if err != nil {
		l.Logger.Error("获取团队已分配算力失败: ", err)
		return nil, err
	}
	return &user.OrgUserPageListResp{
		PageNum:         in.PageNum,
		PageSize:        in.PageSize,
		Total:           total,
		MonthPowerLimit: bsOrg.MonthPowerLimit,
		MonthUsedPower:  bsOrg.MonthUsedPower,
		CanGivePower:    bsOrg.MonthPowerLimit - totalPower,
		List:            orgUserList,
	}, nil
}

func GetUserRolePermList(userId int64, dborgpermission *model.Dborgpermission) ([]*org.UserRole, []*org.UserPerm) {
	var userRoleList []*org.UserRole
	var userPermList []*org.UserPerm
	var orgUser *model.User
	for _, mongoUser := range dborgpermission.UserList {
		if mongoUser.Id == userId {
			orgUser = mongoUser
		}
	}
	if orgUser == nil {
		return userRoleList, userPermList
	}

	// 获取角色列表
	roleMap := make(map[int64]*model.Role)
	for _, mongoRole := range dborgpermission.RoleList {
		roleMap[mongoRole.Id] = mongoRole
	}
	for _, roleId := range orgUser.RoleList {
		role := roleMap[*roleId]
		if role == nil {
			continue
		}
		userRoleList = append(userRoleList, &org.UserRole{
			RoleId:   role.Id,
			RoleName: role.Name,
		})
	}

	// 获取权限列表
	permMap := make(map[int64]*model.Permission)
	for _, mongoPerm := range dborgpermission.PermissionList {
		permMap[mongoPerm.Id] = mongoPerm
	}
	// 先使用map保存防止重复
	uniquePerms := make(map[int64]*org.UserPerm)
	for _, userRole := range userRoleList {
		role := roleMap[userRole.RoleId]
		if role == nil {
			continue
		}
		for _, permId := range role.PermissionList {
			perm := permMap[*permId]
			uniquePerms[*permId] = &org.UserPerm{
				PermId:   perm.Id,
				PermName: perm.Name,
				PermCode: perm.Perm,
			}
		}
	}
	// 转换为切片
	for _, perm := range uniquePerms {
		userPermList = append(userPermList, perm)
	}

	return userRoleList, userPermList
}
