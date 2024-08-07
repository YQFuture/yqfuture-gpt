package org

import (
	"context"
	"encoding/json"
	"strconv"
	"yufuture-gpt/app/user/cmd/rpc/client/org"
	"yufuture-gpt/app/user/model/redis"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrgShopPageListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetOrgShopPageListLogic 获取团队店铺分页列表
func NewGetOrgShopPageListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrgShopPageListLogic {
	return &GetOrgShopPageListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrgShopPageListLogic) GetOrgShopPageList(req *types.OrgShopPageListReq) (resp *types.OrgShopPageListResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.OrgShopPageListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 调用RPC接口 获取团队店铺分页列表
	shopPageListResp, err := l.svcCtx.OrgClient.GetOrgShopPageList(l.ctx, &org.OrgShopPageListReq{
		UserId:    userId,
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		Query:     req.Query,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})
	if err != nil {
		l.Logger.Error("获取团队店铺分页列表失败", err)
		return &types.OrgShopPageListResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "操作失败",
			},
		}, nil
	}

	// 解析构建返回数据
	orgShopPageListResp := &types.OrgShopPageListResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "操作成功",
		},
		Data: types.OrgShopPage{
			BasePageResp: types.BasePageResp{
				PageNum:  shopPageListResp.PageNum,
				PageSize: shopPageListResp.PageSize,
				Total:    shopPageListResp.Total,
			},
			PlatformTypeNum: shopPageListResp.PlatformTypeNum,
		},
	}

	var shopList []types.OrgShop
	for _, shop := range shopPageListResp.List {
		orgShop := types.OrgShop{
			ShopId:          strconv.FormatInt(shop.ShopId, 10),
			PlatformType:    shop.PlatformType,
			ShopName:        shop.ShopName,
			MonthPowerLimit: shop.MonthPowerLimit,
			MonthUsedPower:  shop.MonthUsedPower,
			UserNum:         shop.UserNum,
			CareTime:        shop.CareTime,
			CareTimes:       shop.CareTimes,
			UsedPower:       shop.UsedPower,
		}
		// 角色列表
		var roleList []types.ShopRole
		for _, role := range shop.RoleList {
			shopRole := types.ShopRole{
				RoleId:   strconv.FormatInt(role.RoleId, 10),
				RoleName: role.RoleName,
			}
			roleList = append(roleList, shopRole)
		}
		orgShop.RoleList = roleList

		// 在线用户列表
		var onlineUserList []types.ShopUser
		for _, user := range shop.RoleUserList {
			// 从Redis中获取当前登录用户数据
			loginUser, err := redis.GetLoginUser(l.ctx, l.svcCtx.Redis, strconv.FormatInt(userId, 10))
			if err != nil {
				l.Logger.Error("获取当前登录用户数据失败", err)
				break
			}
			if loginUser == "" {
				break
			}
			shopUser := types.ShopUser{
				UserId:   strconv.FormatInt(user.UserId, 10),
				NickName: user.NickName,
				Phone:    user.Phone,
				HeadImg:  user.HeadImg,
			}
			onlineUserList = append(onlineUserList, shopUser)
		}
		orgShop.OnlineUserList = onlineUserList

		// 关键词转接用户列表
		var keywordSwitchingUserList []types.ShopUser
		for _, user := range shop.KeywordSwitchingUserList {
			shopUser := types.ShopUser{
				UserId:   strconv.FormatInt(user.UserId, 10),
				NickName: user.NickName,
				Phone:    user.Phone,
				HeadImg:  user.HeadImg,
			}
			keywordSwitchingUserList = append(keywordSwitchingUserList, shopUser)
		}
		orgShop.KeywordSwitchingUserList = keywordSwitchingUserList

		// 异常责任人用户列表
		var exceptionDutyUserList []types.ShopUser
		for _, user := range shop.ExceptionDutyUserList {
			shopUser := types.ShopUser{
				UserId:   strconv.FormatInt(user.UserId, 10),
				NickName: user.NickName,
				Phone:    user.Phone,
				HeadImg:  user.HeadImg,
			}
			exceptionDutyUserList = append(exceptionDutyUserList, shopUser)
		}
		orgShop.ExceptionDutyUserList = exceptionDutyUserList

		shopList = append(shopList, orgShop)
	}

	orgShopPageListResp.Data.List = shopList
	return orgShopPageListResp, nil
}
