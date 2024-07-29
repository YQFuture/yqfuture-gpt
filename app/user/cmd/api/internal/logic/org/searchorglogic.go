package org

import (
	"context"
	"strconv"
	"yufuture-gpt/app/user/cmd/rpc/client/org"
	"yufuture-gpt/common/consts"

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
	orgResp, err := l.svcCtx.OrgClient.SearchOrg(l.ctx, &org.SearchOrgReq{
		Query: req.Query,
	})
	if err != nil {
		l.Logger.Error("查找机构失败", err)
		return &types.SearchOrgResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "查找失败",
			},
		}, nil
	}

	var orgInfoList []types.SearchOrgInfo
	for _, orgInfo := range orgResp.Result {
		orgInfoList = append(orgInfoList, types.SearchOrgInfo{
			OrgId:   strconv.FormatInt(orgInfo.OrgId, 10),
			OrgName: orgInfo.OrgName,
		})
	}
	return &types.SearchOrgResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "查找成功·",
		},
		Data: orgInfoList,
	}, nil
}
