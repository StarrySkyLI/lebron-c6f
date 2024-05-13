package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zhoushuguang/lebron/apps/seckill/api/internal/logic"
	"github.com/zhoushuguang/lebron/apps/seckill/api/internal/svc"
	"github.com/zhoushuguang/lebron/apps/seckill/api/internal/types"
)

func QuerResultHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QuerReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewQuerResultLogic(r.Context(), svcCtx)
		resp, err := l.QuerResult(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
