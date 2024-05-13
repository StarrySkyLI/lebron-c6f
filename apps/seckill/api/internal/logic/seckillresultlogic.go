package logic

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/zhoushuguang/lebron/apps/order/rpc/order"
	"github.com/zhoushuguang/lebron/apps/product/rpc/product"
	"github.com/zhoushuguang/lebron/pkg/xerr"

	"github.com/zhoushuguang/lebron/apps/seckill/api/internal/svc"
	"github.com/zhoushuguang/lebron/apps/seckill/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SeckillResultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSeckillResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SeckillResultLogic {
	return &SeckillResultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SeckillResultLogic) SeckillResult(req *types.SeckillResultReq) (resp *types.SeckillResultResponse, err error) {
	key := fmt.Sprintf("skillToken:productId_%d_userid_%d", req.Product_id, req.User_id)
	token, err := l.svcCtx.BizRedis.Get(key)
	if err != nil {
		return nil, err
	}
	if req.SeckillToken != token {
		return nil, errors.Wrap(xerr.NewErrMsg("秒杀令牌校验失败"), "秒杀令牌校验失败")
	}

	_, err = l.svcCtx.OrderRPC.CreateOrder(context.Background(), &order.CreateOrderRequest{Uid: req.User_id, Pid: req.Product_id})
	if err != nil {
		logx.Errorf("CreateOrder uid: %d pid: %d error: %v", req.User_id, req.Product_id, err)
		return nil, err
	}
	_, err = l.svcCtx.ProductRPC.UpdateProductStock(context.Background(), &product.UpdateProductStockRequest{ProductId: req.Product_id, Num: 1})
	if err != nil {
		logx.Errorf("UpdateProductStock uid: %d pid: %d error: %v", req.User_id, req.Product_id, err)
		return nil, err
	}
	_, err = l.svcCtx.BizRedis.Del(key)
	if err != nil {
		return nil, err
	}
	return &types.SeckillResultResponse{}, nil
}
