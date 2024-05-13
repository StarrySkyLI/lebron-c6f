package logic

import (
	"context"
	"github.com/zhoushuguang/lebron/apps/seckill/rpc/seckill"
	"google.golang.org/grpc/status"

	"github.com/zhoushuguang/lebron/apps/seckill/api/internal/svc"
	"github.com/zhoushuguang/lebron/apps/seckill/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SeckillLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSeckillLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillLogic {
	return &SeckillLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SeckillLogic) Seckill(req *types.SeckillReq) (resp *types.SeckillResponse, err error) {
	//id为空报错
	if req.User_id == 0 {
		return nil, status.Error(801, "user_id为空")
	}
	if req.Product_id == 0 {
		return nil, status.Error(801, "product_id为空")
	}
	_, err = l.svcCtx.SeckillRPC.SeckillOrder(l.ctx, &seckill.SeckillOrderRequest{
		UserId:    req.User_id,
		ProductId: req.Product_id,
	})
	if err != nil {
		logx.Errorf("l.svcCtx.SeckillRpc.SeckillOrder error :%v", err)
		return
	}
	return &types.SeckillResponse{}, nil
}
