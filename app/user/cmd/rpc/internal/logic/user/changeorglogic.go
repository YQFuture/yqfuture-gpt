package userlogic

import (
	"context"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeOrgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeOrgLogic {
	return &ChangeOrgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ChangeOrg 切换组织
func (l *ChangeOrgLogic) ChangeOrg(in *user.ChangeOrgReq) (*user.ChangeOrgResp, error) {
	// 判断用户是否在组织中
	org, err := l.svcCtx.BsOrganizationModel.FindOneByIdAndUserId(l.ctx, in.OrgId, in.UserId)
	if err != nil {
		l.Logger.Error("根据组织ID和用户ID获取组织失败", err)
		return nil, err
	}
	if org == nil {
		return &user.ChangeOrgResp{
			Code: consts.UserNotInOrg,
		}, nil
	}

	// 修改用户表中当前组织ID
	err = l.svcCtx.BsUserModel.ChangeOrg(l.ctx, in.OrgId, in.UserId)
	if err != nil {
		l.Logger.Error("修改用户表中当前组织ID失败", err)
		return nil, err
	}

	return &user.ChangeOrgResp{
		Code: consts.Success,
	}, nil
}
