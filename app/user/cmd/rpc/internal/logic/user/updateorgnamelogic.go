package userlogic

import (
	"context"
	"errors"
	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateOrgNameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateOrgNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOrgNameLogic {
	return &UpdateOrgNameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateOrgName 更新组织名称
func (l *UpdateOrgNameLogic) UpdateOrgName(in *user.UpdateOrgNameReq) (*user.UpdateOrgNameResp, error) {
	// 只允许管理员修改组织名称
	org, err := l.svcCtx.BsOrganizationModel.FindOne(l.ctx, in.OrgId)
	if err != nil {
		return nil, err
	}
	if org.OwnerId != in.UserId {
		return nil, errors.New("只允许管理员修改组织名称")
	}

	// 更新组织名称
	err = l.svcCtx.BsOrganizationModel.UpdateOrgName(l.ctx, in.OrgName, in.OrgId)
	if err != nil {
		l.Logger.Error("更新组织名称失败: ", err)
		return nil, err
	}
	return &user.UpdateOrgNameResp{}, nil
}
