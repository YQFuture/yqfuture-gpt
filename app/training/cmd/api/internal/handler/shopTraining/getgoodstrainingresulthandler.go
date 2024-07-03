package shopTraining

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yufuture-gpt/app/training/cmd/api/internal/logic/shopTraining"
	"yufuture-gpt/app/training/cmd/api/internal/svc"
	"yufuture-gpt/app/training/cmd/api/internal/types"
)

// 获取商品训练结果
func GetGoodsTrainingResultHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetGoodsTrainingResultReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := shopTraining.NewGetGoodsTrainingResultLogic(r.Context(), svcCtx)
		resp, err := l.GetGoodsTrainingResult(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
