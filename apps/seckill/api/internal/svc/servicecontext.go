package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zhoushuguang/lebron/apps/order/rpc/order"
	"github.com/zhoushuguang/lebron/apps/product/rpc/product"
	"github.com/zhoushuguang/lebron/apps/seckill/api/internal/config"
	"github.com/zhoushuguang/lebron/apps/seckill/rpc/seckill"
)

type ServiceContext struct {
	Config     config.Config
	SeckillRPC seckill.Seckill
	ProductRPC product.Product
	OrderRPC   order.Order
	BizRedis   *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {

	return &ServiceContext{
		Config:     c,
		SeckillRPC: seckill.NewSeckill(zrpc.MustNewClient(c.SeckillRPC)),
		ProductRPC: product.NewProduct(zrpc.MustNewClient(c.ProductRPC)),
		OrderRPC:   order.NewOrder(zrpc.MustNewClient(c.OrderRPC)),
		BizRedis:   redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
	}
}
