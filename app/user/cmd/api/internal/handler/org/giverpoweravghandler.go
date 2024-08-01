package org

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yufuture-gpt/app/user/cmd/api/internal/logic/org"
	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"
)

// 平均分配算力
func GiverPowerAvgHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GivePowerAvgReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := org.NewGiverPowerAvgLogic(r.Context(), svcCtx)
		resp, err := l.GiverPowerAvg(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
