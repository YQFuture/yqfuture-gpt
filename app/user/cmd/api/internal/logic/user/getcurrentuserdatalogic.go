package user

import (
	"context"
	"encoding/json"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"
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
		l.Logger.Error("获取用户id失败", err)
		return &types.CurrentUserDataResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "获取失败",
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
	if currentUserDataResp.Code == consts.PhoneTsNotBound {
		return &types.CurrentUserDataResp{
			BaseResp: types.BaseResp{
				Code: consts.PhoneTsNotBound,
				Msg:  "手机号未绑定",
			},
		}, nil
	}

	return &types.CurrentUserDataResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "获取成功",
		},
		Data: types.CurrentUserData{},
	}, nil
}
