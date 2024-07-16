package login

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yufuture-gpt/app/user/cmd/api/internal/logic/login"
	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"
)

// 获取图形验证码
func GetCaptchaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BaseReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := login.NewGetCaptchaLogic(r.Context(), svcCtx)
		resp, err := l.GetCaptcha(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
