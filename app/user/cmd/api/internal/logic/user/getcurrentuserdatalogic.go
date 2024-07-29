package user

import (
	"context"
	"encoding/json"
	"strconv"
	"time"
	"yufuture-gpt/app/user/cmd/api/internal/logic/login"
	"yufuture-gpt/app/user/cmd/rpc/client/user"
	"yufuture-gpt/app/user/model/redis"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCurrentUserDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetCurrentUserDataLogic 获取当前登录用户数据
func NewGetCurrentUserDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserDataLogic {
	return &GetCurrentUserDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCurrentUserDataLogic) GetCurrentUserData(req *types.BaseReq) (resp *types.CurrentUserDataResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.CurrentUserDataResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "获取失败",
			},
		}, nil
	}

	// 获取当前登录用户数据
	loginUser, err := redis.GetLoginUser(l.ctx, l.svcCtx.Redis, strconv.FormatInt(userId, 10))
	if err != nil {
		return nil, err
	}
	if loginUser == "" {
		return &types.CurrentUserDataResp{
			BaseResp: types.BaseResp{
				Code: consts.Unauthorized,
				Msg:  "登录失效",
			},
		}, nil
	}

	// 调用RPC接口 获取当前登录用户数据
	currentUserDataResp, err := l.svcCtx.UserClient.GetCurrentUserData(l.ctx, &user.CurrentUserDataReq{
		UserId: userId,
	})
	if err != nil {
		l.Logger.Error("获取当前登录用户数据失败", err)
		return &types.CurrentUserDataResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "获取失败",
			},
		}, nil
	}
	// 判断是否绑定手机号
	if currentUserDataResp.Code == consts.PhoneIsNotBound {
		return &types.CurrentUserDataResp{
			BaseResp: types.BaseResp{
				Code: consts.PhoneIsNotBound,
				Msg:  "手机号未绑定",
			},
		}, nil
	}

	// 生成 Token
	accessExpire := l.svcCtx.Config.Auth.AccessExpire
	payload := map[string]interface{}{
		"id":      currentUserDataResp.Result.Id,
		"ex_time": time.Now().AddDate(0, 0, 7),
	}
	token, err := login.GetJwtToken(l.svcCtx.Config.Auth.AccessSecret, accessExpire, payload)
	if err != nil {
		l.Logger.Error("生成token失败", err)
		return &types.CurrentUserDataResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "获取失败 请重试",
			},
		}, nil
	}

	return &types.CurrentUserDataResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "获取成功",
		},
		Data: types.CurrentUserData{
			Token:    token,
			Phone:    currentUserDataResp.Result.Phone,
			NickName: currentUserDataResp.Result.NickName,
			HeadImg:  currentUserDataResp.Result.HeadImg,
			NowOrg: types.OrgInfo{
				OrgId:      strconv.FormatInt(currentUserDataResp.Result.NowOrg.OrgId, 10),
				OrgName:    currentUserDataResp.Result.NowOrg.OrgName,
				BundleType: currentUserDataResp.Result.NowOrg.BundleType,
				IsAdmin:    currentUserDataResp.Result.NowOrg.IsAdmin,
			},
			UnreadMsgCount: currentUserDataResp.Result.UnreadMsgCount,
		},
	}, nil
}
