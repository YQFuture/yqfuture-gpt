package orglogic

import (
	"context"
	"errors"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrgPermListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrgPermListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgPermListLogic {
	return &GetOrgPermListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetOrgPermList 获取团队权限列表
func (l *GetOrgPermListLogic) GetOrgPermList(in *user.OrgPermListReq) (*user.OrgPermListResp, error) {
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
	permList := dborgpermission.PermissionList

	// 封装返回数据
	var orgPermList []*user.OrgPerm
	for _, perm := range permList {
		orgPermList = append(orgPermList, &user.OrgPerm{
			PermId:   perm.Id,
			PermName: perm.Name,
			PermCode: perm.Perm,
			ParentId: perm.ParentId,
			Sort:     perm.Sort,
		})
	}

	return &user.OrgPermListResp{
		Result: orgPermList,
	}, nil
}
