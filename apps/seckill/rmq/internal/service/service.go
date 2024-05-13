package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"github.com/zhoushuguang/lebron/apps/order/rpc/order"
	"github.com/zhoushuguang/lebron/apps/product/rpc/product"
	"github.com/zhoushuguang/lebron/apps/seckill/rmq/internal/config"
	"log"
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	chanCount   = 10
	bufferCount = 1024
)

type Service struct {
	c          config.Config
	ProductRPC product.Product
	OrderRPC   order.Order
	BizRedis   *redis.Redis
	waiter     sync.WaitGroup
	msgsChan   []chan *KafkaData
}

type KafkaData struct {
	Uid int64 `json:"uid"`
	Pid int64 `json:"pid"`
}

func NewService(c config.Config) *Service {
	s := &Service{
		c:          c,
		BizRedis:   redis.New(c.BizRedis.Host, redis.WithPass(c.BizRedis.Pass)),
		ProductRPC: product.NewProduct(zrpc.MustNewClient(c.ProductRPC)),
		OrderRPC:   order.NewOrder(zrpc.MustNewClient(c.OrderRPC)),
		msgsChan:   make([]chan *KafkaData, chanCount),
	}
	for i := 0; i < chanCount; i++ {
		ch := make(chan *KafkaData, bufferCount)
		s.msgsChan[i] = ch
		s.waiter.Add(1)
		go s.consume(ch)
	}

	return s
}

func (s *Service) consume(ch chan *KafkaData) {
	defer s.waiter.Done()

	for {
		m, ok := <-ch
		if !ok {
			log.Fatal("seckill rmq exit")
		}
		fmt.Printf("consume msg: %+v\n", m)
		_, err := s.ProductRPC.CheckAndUpdateStock(context.Background(), &product.CheckAndUpdateStockRequest{ProductId: m.Pid})
		if err != nil {
			logx.Errorf("s.ProductRPC.CheckAndUpdateStock pid: %d error: %v", m.Pid, err)

			continue
		}
		//	//生成秒杀token
		_ = s.generateSecondKillToken(m.Pid, m.Uid)

	}
}

func (s *Service) Consume(_ string, value string) error {
	logx.Infof("Consume value: %s\n", value)
	var data []*KafkaData
	if err := json.Unmarshal([]byte(value), &data); err != nil {
		return err
	}
	for _, d := range data {
		s.msgsChan[d.Pid%chanCount] <- d
	}
	return nil
}

func (s *Service) generateSecondKillToken(productId, userId int64) error {
	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}
	key := fmt.Sprintf("skillToken:productId_%d_userid_%d", productId, userId)
	err = s.BizRedis.Set(key, generateUUID)
	if err != nil {
		return err
	}
	err = s.BizRedis.Expire(key, 5*60)
	if err != nil {
		return err
	}
	return nil
}
