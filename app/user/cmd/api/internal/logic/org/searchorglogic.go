package org

import (
	"context"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SearchOrgLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewSearchOrgLogic 查找团队
func NewSearchOrgLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchOrgLogic {
	return &SearchOrgLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchOrgLogic) SearchOrg(req *types.SearchOrgReq) (resp *types.SearchOrgResp, err error) {
	// todo: add your logic here and delete this line

	return
}
