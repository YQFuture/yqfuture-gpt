package orglogic

import (
	"context"
	"errors"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserLogic {
	return &DeleteUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DeleteUser 删除用户
func (l *DeleteUserLogic) DeleteUser(in *user.DeleteUserReq) (*user.DeleteUserResp, error) {
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
		return nil, errors.New("只有团队管理员才能删除用户")
	}
	// 调用MongoDB获取团队权限文档
	dborgpermission, err := l.svcCtx.DborgpermissionModel.FindOne(l.ctx, bsOrg.MongoPermId)
	if err != nil {
		l.Logger.Error("获取团队权限文档失败: ", err)
		return nil, err
	}

	// 删除MongoDB文档中的用户
	for index, mongoUser := range dborgpermission.UserList {
		if mongoUser.Id == in.DeleteUserId {
			dborgpermission.UserList = append(dborgpermission.UserList[:index], dborgpermission.UserList[index+1:]...)
			break
		}
	}
	// 更新MongoDB中的团队权限文档
	_, err = l.svcCtx.DborgpermissionModel.Update(l.ctx, dborgpermission)
	if err != nil {
		l.Logger.Error("更新团队权限文档失败: ", err)
		return nil, err
	}

	// 删除用户组织关联表中的用户
	err = l.svcCtx.BsUserOrgModel.DeleteByUserIdAndOrgId(l.ctx, in.DeleteUserId, bsOrg.Id)
	if err != nil {
		l.Logger.Error("删除用户组织关联表中的用户失败: ", err)
		return nil, err
	}

	return &user.DeleteUserResp{}, nil
}
