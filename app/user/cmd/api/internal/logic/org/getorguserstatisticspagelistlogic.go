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

type GetOrgUserStatisticsPageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetOrgUserStatisticsPageListLogic 获取团队用户统计信息分页列表
func NewGetOrgUserStatisticsPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgUserStatisticsPageListLogic {
	return &GetOrgUserStatisticsPageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrgUserStatisticsPageListLogic) GetOrgUserStatisticsPageList(req *types.OrgUserStatisticsPageListReq) (resp *types.OrgUserStatisticsPageListResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.OrgUserStatisticsPageListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}
	shopId, err := strconv.ParseInt(req.ShopId, 10, 64)
	if err != nil {
		l.Logger.Error("获取店铺ID失败", err)
		return &types.OrgUserStatisticsPageListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 调用RPC接口 获取团队用户统计信息分页列表
	pageListResp, err := l.svcCtx.OrgClient.GetOrgUserStatisticsPageList(l.ctx, &org.OrgUserStatisticsPageListReq{
		UserId:    userId,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		ShopId:    shopId,
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		Query:     req.Query,
	})
	if err != nil {
		l.Logger.Error("获取团队用户统计信息分页列表失败", err)
		return &types.OrgUserStatisticsPageListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 封装返回数据
	resp = &types.OrgUserStatisticsPageListResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
		Data: types.OrgUserStatisticsPage{
			BasePageResp: types.BasePageResp{
				PageNum:  pageListResp.PageNum,
				PageSize: pageListResp.PageSize,
				Total:    pageListResp.Total,
			},
			List: make([]types.OrgUserStatistics, 0),
		},
	}
	var userStatisticsList []types.OrgUserStatistics
	for _, item := range pageListResp.List {
		userStatistics := types.OrgUserStatistics{
			UserId:           strconv.FormatInt(item.UserId, 10),
			Phone:            item.Phone,
			NickName:         item.NickName,
			HeadImg:          item.HeadImg,
			CareTime:         item.CareTime,
			CareTimes:        item.CareTimes,
			UsedPower:        item.UsedPower,
			RecentOnlineTime: item.RecentOnlineTime,
			TotalOnlineTime:  item.TotalOnlineTime,
		}
		var shopList []types.RoleShop
		for _, shop := range item.ShopList {
			shopList = append(shopList, types.RoleShop{
				ShopId:       strconv.FormatInt(shop.ShopId, 10),
				ShopName:     shop.ShopName,
				PlatformType: shop.PlatformType,
			})
		}
		userStatistics.ShopList = shopList

		userStatisticsList = append(userStatisticsList, userStatistics)
	}
	resp.Data.List = userStatisticsList

	return resp, nil
}
