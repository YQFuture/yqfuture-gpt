package login

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"sort"
	"strings"

	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type WechatCallBackGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewWechatCallBackGetLogic 微信回调
func NewWechatCallBackGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WechatCallBackGetLogic {
	return &WechatCallBackGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WechatCallBackGetLogic) WechatCallBackGet(req *types.WechatCallBackGetReq) (resp string, err error) {
	// 验证消息来源
	if !CheckSignature(l.svcCtx.Config.WechatConf.Token, req.Signature, req.Timestamp, req.Nonce) {
		return "参数错误", nil
	}
	return req.Echostr, nil
}

func CheckSignature(token, signature, timestamp, nonce string) bool {
	stringList := []string{token, timestamp, nonce}
	// 字典序排序
	sort.Strings(stringList)
	sortedStr := strings.Join(stringList, "")
	h := sha1.New()
	h.Write([]byte(sortedStr))
	calcSignature := hex.EncodeToString(h.Sum(nil))
	return calcSignature == signature
}
