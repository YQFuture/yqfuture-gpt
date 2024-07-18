package login

import (
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
	"yufuture-gpt/app/user/cmd/api/internal/logic/login"
	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"
)

// 微信回调
func WechatCallBackGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.WechatCallBackGetReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := login.NewWechatCallBackGetLogic(r.Context(), svcCtx)
		resp, err := l.WechatCallBackGet(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			//httpx.OkJsonCtx(r.Context(), w, resp)
			w.Header().Set("Content-Type", "text/plain; charset=utf-8")
			_, err := w.Write([]byte(resp))
			if err != nil {
				httpx.ErrorCtx(r.Context(), w, err)
			}
		}
	}
}
