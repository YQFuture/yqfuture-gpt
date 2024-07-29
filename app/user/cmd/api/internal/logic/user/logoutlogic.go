package user

import (
	"context"
	"encoding/json"
	"strconv"
	"yufuture-gpt/app/user/model/redis"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogOutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewLogOutLogic 退出登录
func NewLogOutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogOutLogic {
	return &LogOutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogOutLogic) LogOut(req *types.LogOutReq) (resp *types.LogOutResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户ID失败", err)
		return &types.LogOutResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "退出登录失败",
			},
		}, nil
	}

	err = redis.DelLoginUser(l.ctx, l.svcCtx.Redis, strconv.FormatInt(userId, 10))
	if err != nil {
		l.Logger.Error("删除用户登录状态失败", err)
		return &types.LogOutResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "退出登录失败",
			},
		}, nil
	}

	return &types.LogOutResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "退出登录成功",
		},
	}, nil
}
