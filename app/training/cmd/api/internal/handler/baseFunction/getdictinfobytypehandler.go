package baseFunction

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yufuture-gpt/app/training/cmd/api/internal/logic/baseFunction"
	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"
)

func GetDictInfoByTypeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DictInfoByTypeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := baseFunction.NewGetDictInfoByTypeLogic(r.Context(), svcCtx)
		resp, err := l.GetDictInfoByType(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
