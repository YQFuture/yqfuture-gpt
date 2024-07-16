package login

import (
	"context"
	"github.com/google/uuid"
	"github.com/mojocn/base64Captcha"
	"yufuture-gpt/app/user/model/redis"
	"yufuture-gpt/common/consts"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetCaptchaLogic 获取图形验证码
func NewGetCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCaptchaLogic {
	return &GetCaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCaptchaLogic) GetCaptcha(req *types.BaseReq) (resp *types.CaptchaResp, err error) {
	// 配置验证码参数
	driver := base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)
	// 生成验证码 配置最小的存储容量和失效时间 避免内存占用
	captcha := base64Captcha.NewCaptcha(driver, base64Captcha.NewMemoryStore(1, 1))
	// 生成验证码
	_, b64s, answer, err := captcha.Generate()
	if err != nil {
		l.Logger.Error("生成图形验证码失败: ", err)
		return &types.CaptchaResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "获取图形验证码失败 请重试",
			},
		}, nil
	}

	CaptchaId := uuid.New().String()
	err = redis.SetImgCaptcha(l.ctx, l.svcCtx.Redis, CaptchaId, answer)
	if err != nil {
		l.Logger.Error("保存图形验证码到Redis: ", err)
		return &types.CaptchaResp{
			BaseResp: types.BaseResp{
				Code: consts.Fail,
				Msg:  "获取图形验证码失败 请重试",
			},
		}, nil
	}

	return &types.CaptchaResp{
		BaseResp: types.BaseResp{
			Code: consts.Success,
			Msg:  "获取图形验证码成功",
		},
		Data: types.CaptchaData{
			CaptchaId:  CaptchaId,
			CaptchaImg: b64s,
		},
	}, nil
}
