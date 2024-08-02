package org

import (
	"context"
	"encoding/json"
	"strconv"
	"yufuture-gpt/app/user/cmd/rpc/client/org"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrgUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetOrgUserListLogic 获取团队用户列表
func NewGetOrgUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgUserListLogic {
	return &GetOrgUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrgUserListLogic) GetOrgUserList(req *types.OrgUserListReq) (resp *types.OrgUserListResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.OrgUserListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 调用RPC接口 获取团队用户列表
	orgUserListResp, err := l.svcCtx.OrgClient.GetOrgUserList(l.ctx, &org.OrgUserListReq{
		UserId: userId,
	})
	if err != nil {
		l.Logger.Error("获取团队用户列表失败", err)
		return &types.OrgUserListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	var orgUserList []types.OrgUser
	for _, orgUserResp := range orgUserListResp.List {
		orgUser := types.OrgUser{
			UserId:   strconv.FormatInt(orgUserResp.UserId, 10),
			Phone:    orgUserResp.Phone,
			NickName: orgUserResp.NickName,
			HeadImg:  orgUserResp.HeadImg,
			Status:   orgUserResp.Status,
		}
		orgUserList = append(orgUserList, orgUser)
	}

	return &types.OrgUserListResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
		Data: orgUserList,
	}, nil
}
