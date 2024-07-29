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
	// todo: add your logic here and delete this line

	return &user.SearchOrgReqResp{}, nil
}
