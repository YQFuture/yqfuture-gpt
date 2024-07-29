package orglogic

import (
	"context"

	"yufuture-gpt/app/user/cmd/rpc/internal/svc"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchOrgLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSearchOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchOrgLogic {
	return &SearchOrgLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SearchOrg 查找团队
func (l *SearchOrgLogic) SearchOrg(in *user.SearchOrgReq) (*user.SearchOrgReqResp, error) {
	orgList, err := l.svcCtx.BsOrganizationModel.FindListByNameOrOwnerPhone(l.ctx, in.Query)
	if err != nil {
		l.Logger.Error("根据用户名或手机号查找团队失败", err)
		return nil, err
	}

	var orgInfoList []*user.SearchOrgInfo
	for _, org := range *orgList {
		orgInfoList = append(orgInfoList, &user.SearchOrgInfo{
			OrgId:   org.Id,
			OrgName: org.OrgName.String,
		})
	}
	return &user.SearchOrgReqResp{
		Result: orgInfoList,
	}, nil
}
