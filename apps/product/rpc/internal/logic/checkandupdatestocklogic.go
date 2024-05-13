package logic

import (
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strconv"

	"github.com/zhoushuguang/lebron/apps/product/rpc/internal/svc"
	"github.com/zhoushuguang/lebron/apps/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type CheckAndUpdateStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCheckAndUpdateStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckAndUpdateStockLogic {

	return &CheckAndUpdateStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

const (
	luaCheckAndUpdateScript = `
local counts = redis.call("HMGET", KEYS[1], "totalCount", "seckillCount")
local total = tonumber(counts[1])
local seckill = tonumber(counts[2])
if seckill + 1 <= total then
	redis.call("HINCRBY", KEYS[1], "seckillCount", 1)
	return 1
end
return 0
`
)

// 1. 通过`HMGET`命令获取特定商品(key是`KEYS[1]`)的总库存(`totalCount`)和已秒杀数量(`seckillCount`)。  库存检查
// 2. 将获取到的库存和已秒杀数量转换为数字。
// 3. 判断如果已秒杀数量加1后仍然不超过总库存，则表示还有库存可以进行秒杀。
// 1. 使用`HINCRBY`命令将`seckillCount`(已秒杀数量)增加1。                                          库存扣减
// 2. 返回`1`表示秒杀成功。
// 4. 如果已秒杀数量加1后超过总库存，则没有库存进行秒杀，直接返回`0`表示秒杀失败。
func (l *CheckAndUpdateStockLogic) CheckAndUpdateStock(in *product.CheckAndUpdateStockRequest) (*product.CheckAndUpdateStockResponse, error) {
	//l.prepareData()

	val, err := l.svcCtx.BizRedis.EvalCtx(l.ctx, luaCheckAndUpdateScript, []string{stockKey(in.ProductId)})
	if err != nil {
		return nil, err
	}
	if val.(int64) == 0 {
		return nil, status.Errorf(codes.ResourceExhausted, fmt.Sprintf("insufficient stock: %d", in.ProductId))
	}

	return &product.CheckAndUpdateStockResponse{}, nil
}

func stockKey(pid int64) string {
	return fmt.Sprintf("seckill:goodsStock:%d", pid)
}

type SeckillInfo struct {
	TotalCount   int `json:"totalCount"`
	InitStatus   int `json:"initStatus"`
	SeckillCount int `json:"seckillCount"`
}

func (l *CheckAndUpdateStockLogic) prepareData() {
	goodsID := "10" // 替换为实际的商品 ID
	key := fmt.Sprintf("seckill:goodsStock:%s", goodsID)

	// 创建 SeckillInfo 结构体
	seckillData := SeckillInfo{
		TotalCount:   200,
		InitStatus:   0,
		SeckillCount: 0,
	}

	// 将结构体转换为 map[string]interface{} 格式
	dataMap := map[string]string{
		"totalCount":   strconv.Itoa(seckillData.TotalCount),
		"initStatus":   strconv.Itoa(seckillData.InitStatus),
		"seckillCount": strconv.Itoa(seckillData.SeckillCount),
	}

	err := l.svcCtx.BizRedis.Hmset(key, dataMap)
	if err != nil {
		return
	}

}
