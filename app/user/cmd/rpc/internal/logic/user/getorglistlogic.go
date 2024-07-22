package userlogic

import (
	"context"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrgListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrgListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgListLogic {
	return &GetOrgListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetOrgList 获取用户组织列表
func (l *GetOrgListLogic) GetOrgList(in *user.OrgListReq) (*user.OrgListResp, error) {
	orgList, err := l.svcCtx.BsOrganizationModel.FindListByUserId(l.ctx, in.UserId)
	if err != nil {
		l.Logger.Error("获取用户组织列表失败", err)
		return nil, err
	}
	var result []*user.OrgInfo
	for _, org := range *orgList {
		result = append(result, &user.OrgInfo{
			OrgId:      org.Id,
			OrgName:    org.OrgName.String,
			BundleType: org.BundleType,
			IsAdmin:    org.OwnerId == in.UserId,
		})
	}
	return &user.OrgListResp{
		Result: result,
	}, nil
}
