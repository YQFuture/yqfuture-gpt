package user

import (
	"context"
	"encoding/json"
	"yufuture-gpt/app/user/cmd/rpc/pb/user"
	"yufuture-gpt/app/user/model/redis"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindPhoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewBindPhoneLogic 绑定手机号码
func NewBindPhoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindPhoneLogic {
	return &BindPhoneLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *BindPhoneLogic) BindPhone(req *types.BindPhoneReq) (resp *types.BaseResp, err error) {
	id := l.ctx.Value("id")
	userId, err := id.(json.Number).Int64()
	if err != nil {
		l.Logger.Error("获取用户id失败", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "绑定失败 请重试",
		}, nil
	}

	// 验证手机短信验证码
	code, err := redis.GetVerificationCode(l.ctx, l.svcCtx.Redis, req.Phone)
	if err != nil {
		l.Logger.Error("从Redis中获取手机短信验证码失败", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "绑定失败 请重试",
		}, nil
	}
	if code != req.VerificationCode {
		return &types.BaseResp{
			Code: consts.IncorrectVerificationCode,
			Msg:  "手机短信验证码错误",
		}, nil
	}

	// 调用RPC接口 绑定手机号码
	_, err = l.svcCtx.UserClient.BindPhone(l.ctx, &user.BindPhoneReq{
		UserId: userId,
		Phone:  req.Phone,
	})
	if err != nil {
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "绑定失败 请重试",
		}, nil
	}

	return &types.BaseResp{
		Code: consts.Success,
		Msg:  "绑定成功",
	}, nil
}
