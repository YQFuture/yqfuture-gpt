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

type UpdateShopAssignLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUpdateShopAssignLogic 编辑店铺指派
func NewUpdateShopAssignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateShopAssignLogic {
	return &UpdateShopAssignLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateShopAssignLogic) UpdateShopAssign(req *types.UpdateShopAssignReq) (resp *types.UpdateShopAssignResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.UpdateShopAssignResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	shopId, err := strconv.ParseInt(req.ShopId, 10, 64)
	if err != nil {
		l.Logger.Error("获取店铺ID失败", err)
		return &types.UpdateShopAssignResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	var keywordSwitchingUserList []int64
	for _, v := range req.KeywordSwitchingUserList {
		orgUserId, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			l.Logger.Error("获取用户ID失败", err)
			return &types.UpdateShopAssignResp{
				BaseResp: types.BaseResp{
					Code: consts.Fail,
					Msg:  "操作失败",
				},
			}, nil
		}
		keywordSwitchingUserList = append(keywordSwitchingUserList, orgUserId)
	}
	var exceptionDutyUserList []int64
	for _, v := range req.ExceptionDutyUserList {
		orgUserId, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			l.Logger.Error("获取用户ID失败", err)
			return &types.UpdateShopAssignResp{
				BaseResp: types.BaseResp{
					Code: consts.Fail,
					Msg:  "操作失败",
				},
			}, nil
		}
		exceptionDutyUserList = append(exceptionDutyUserList, orgUserId)
	}
	var roleList []int64
	for _, v := range req.RoleList {
		roleId, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			l.Logger.Error("获取角色ID失败", err)
			return &types.UpdateShopAssignResp{
				BaseResp: types.BaseResp{
					Code: consts.Fail,
					Msg:  "操作失败",
				},
			}, nil
		}
		roleList = append(roleList, roleId)
	}

	// 调用RPC接口 编辑店铺指派
	_, err = l.svcCtx.OrgClient.UpdateShopAssign(l.ctx, &org.UpdateShopAssignReq{
		UserId:                   userId,
		ShopId:                   shopId,
		KeywordSwitchingUserList: keywordSwitchingUserList,
		ExceptionDutyUserList:    exceptionDutyUserList,
		RoleList:                 roleList,
	})
	if err != nil {
		l.Logger.Error("编辑店铺指派失败", err)
		return &types.UpdateShopAssignResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	return &types.UpdateShopAssignResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
	}, nil
}
