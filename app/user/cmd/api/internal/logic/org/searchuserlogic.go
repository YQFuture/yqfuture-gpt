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

type SearchUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewSearchUserLogic 查找用户
func NewSearchUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SearchUserLogic {
	return &SearchUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SearchUserLogic) SearchUser(req *types.SearchUserReq) (resp *types.SearchUserResp, err error) {
	userResp, err := l.svcCtx.OrgClient.SearchUser(l.ctx, &org.SearchUserReq{
		Query: req.Query,
	})
	if err != nil {
		l.Logger.Error("查找用户失败", err)
		return &types.SearchUserResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "查找失败",
			},
		}, nil
	}

	var userInfoList []types.SearchUserInfo
	for _, userInfo := range userResp.Result {
		userInfoList = append(userInfoList, types.SearchUserInfo{
			UserId: strconv.FormatInt(userInfo.UserId, 10),
			Phone:  userInfo.Phone,
		})
	}

	return &types.SearchUserResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "查找成功",
		},
		Data: userInfoList,
	}, nil
}
