package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zhoushuguang/lebron/apps/seckill/api/internal/logic"
	"github.com/zhoushuguang/lebron/apps/seckill/api/internal/svc"
	"github.com/zhoushuguang/lebron/apps/seckill/api/internal/types"
)

func SeckillHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SeckillReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.Error(w, err)
			return
		}

		l := logic.NewSeckillLogic(r.Context(), svcCtx)
		resp, err := l.Seckill(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
