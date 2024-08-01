package org

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yufuture-gpt/app/user/cmd/api/internal/logic/org"
	"yufuture-gpt/app/user/cmd/api/internal/svc"
	"yufuture-gpt/app/user/cmd/api/internal/types"
)

// 创建角色
func CreateRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateRoleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := org.NewCreateRoleLogic(r.Context(), svcCtx)
		resp, err := l.CreateRole(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
