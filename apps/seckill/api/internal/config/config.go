package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	SeckillRPC zrpc.RpcClientConf
	ProductRPC zrpc.RpcClientConf
	OrderRPC   zrpc.RpcClientConf
	BizRedis   redis.RedisConf
}
