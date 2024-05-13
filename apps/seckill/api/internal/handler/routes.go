// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	"github.com/zhoushuguang/lebron/apps/seckill/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/sekill",
				Handler: SeckillHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/sekillResult",
				Handler: SeckillResultHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/querresult",
				Handler: QuerResultHandler(serverCtx),
			},
		},
		rest.WithPrefix("/v1"),
	)
}