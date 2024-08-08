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

type GetShopUserListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetShopUserListLogic 获取店铺客服列表
func NewGetShopUserListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetShopUserListLogic {
	return &GetShopUserListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetShopUserListLogic) GetShopUserList(req *types.ShopUserListReq) (resp *types.ShopUserListResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.ShopUserListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	shopId, err := strconv.ParseInt(req.ShopId, 10, 64)
	if err != nil {
		l.Logger.Error("获取店铺ID失败", err)
		return &types.ShopUserListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 调用RPC接口 获取店铺客服列表
	listResp, err := l.svcCtx.OrgClient.GetShopUserList(l.ctx, &org.ShopUserListReq{
		UserId:    userId,
		ShopId:    shopId,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Query:     req.Query,
	})
	if err != nil {
		l.Logger.Error("获取店铺客服列表失败", err)
		return &types.ShopUserListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	var shopUserList []types.ShopUser
	for _, v := range listResp.List {
		shopUserList = append(shopUserList, types.ShopUser{
			UserId:    strconv.FormatInt(v.UserId, 10),
			Phone:     v.Phone,
			NickName:  v.NickName,
			HeadImg:   v.HeadImg,
			CareTime:  v.CareTime,
			CareTimes: v.CareTimes,
			UsedPower: v.UsedPower,
		})
	}

	return &types.ShopUserListResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
		Data: shopUserList,
	}, nil
}
