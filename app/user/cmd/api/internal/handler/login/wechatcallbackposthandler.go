package login

import (
	"encoding/xml"
	"github.com/zeromicro/go-zero/core/logx"
	"io"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yufuture-gpt/app/user/cmd/api/internal/logic/login"
	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"
)

// 微信回调
func WechatCallBackPostHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var signReq types.WechatCallBackPostSignReq
		var req types.WechatCallBackPostReq
		// 获取请求路径中的签名字段
		if err := httpx.Parse(r, &signReq); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		// 反序列化请求体的XML
		body, err := io.ReadAll(r.Body)
		if err != nil {
			logx.Error("读取Body失败 ", err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		err = xml.Unmarshal(body, &req)
		if err != nil {
			logx.Error("反序列化XML失败 ", string(body), err)
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		// 签名字段赋值
		req.Signature = signReq.Signature
		req.Timestamp = signReq.Timestamp
		req.Nonce = signReq.Nonce
		l := login.NewWechatCallBackPostLogic(r.Context(), svcCtx)
		resp, err := l.WechatCallBackPost(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
