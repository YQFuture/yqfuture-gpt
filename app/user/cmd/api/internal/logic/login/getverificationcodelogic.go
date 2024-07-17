package login

import (
	"context"
	"fmt"
	"github.com/alibabacloud-go/tea/tea"
	"math/rand"
	"regexp"
	"time"
	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/thirdparty"
	"yufuture-gpt/app/user/cmd/api/internal/types"
	"yufuture-gpt/app/user/model/redis"
	"yufuture-gpt/common/consts"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v4/client"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetVerificationCodeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewGetVerificationCodeLogic 获取手机短信验证码
func NewGetVerificationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVerificationCodeLogic {
	return &GetVerificationCodeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetVerificationCodeLogic) GetVerificationCode(req *types.VerificationCodeReq) (resp *types.BaseResp, err error) {
	// 验证手机号格式
	regex := `^1[3456789]\d{9}$`
	match, err := regexp.MatchString(regex, req.Phone)
	if err != nil {
		l.Logger.Error("验证手机号格式失败", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "发送失败",
		}, nil
	}
	if !match {
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "手机号格式不正确",
		}, nil
	}

	// 判断是否在允许重试的时间内
	code, err := redis.GetVerificationCode(l.ctx, l.svcCtx.Redis, req.Phone)
	if err != nil {
		l.Logger.Error("从Redis中获取手机短信验证码", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "发送失败",
		}, nil
	}
	if code != "" {
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "请在一分钟后重试",
		}, nil
	}

	// 生成六位短信验证码
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomNumber := rng.Intn(1000000)
	verificationCode := fmt.Sprintf("%06d", randomNumber)

	// 发送短信
	err = thirdparty.SendVerificationCode(l.Logger,
		l.svcCtx.Config.SmsConf.AccessKeyId,
		l.svcCtx.Config.SmsConf.AccessKeySecret,
		l.svcCtx.Config.SmsConf.Domain,
		req.Phone,
		l.svcCtx.Config.SmsConf.SignName,
		l.svcCtx.Config.SmsConf.TemplateCode,
		verificationCode)
	if err != nil {
		l.Logger.Error("发送手机短信失败: ", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "发送失败",
		}, nil
	}

	// 保存手机短信验证码到Redis
	err = redis.SetVerificationCode(l.ctx, l.svcCtx.Redis, req.Phone, verificationCode)
	if err != nil {
		l.Logger.Error("保存手机短信验证码到Redis失败: ", err)
		return &types.BaseResp{
			Code: consts.Fail,
			Msg:  "发送失败",
		}, nil
	}

	return &types.BaseResp{
		Code: consts.Success,
		Msg:  "发送成功",
	}, nil
}

func CreateClient(accessKeyId, accessKeySecret, domain string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     tea.String(accessKeyId),
		AccessKeySecret: tea.String(accessKeySecret),
	}
	config.Endpoint = tea.String(domain)
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}
