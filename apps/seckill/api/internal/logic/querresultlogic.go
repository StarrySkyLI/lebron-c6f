package logic

import (
	"context"
	"fmt"
	"github.com/zhoushuguang/lebron/apps/seckill/api/internal/svc"
	"github.com/zhoushuguang/lebron/apps/seckill/api/internal/types"
	"strconv"

	"github.com/zeromicro/go-zero/core/logx"
)

type QuerResultLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQuerResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuerResultLogic {
	return &QuerResultLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuerResultLogic) QuerResult(req *types.QuerReq) (resp *types.QuerResponse, err error) {
	key := fmt.Sprintf("skillToken:productId_%d_userid_%d", req.Product_id, req.User_id)
	exists, err := l.svcCtx.BizRedis.Exists(key)

	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	//秒杀成功，跳转确认界面
	if exists == true {
		val, err := l.svcCtx.BizRedis.Get(key)
		if err != nil {
			return nil, err
		}
		return &types.QuerResponse{
			Message: val,
		}, nil
	} else {
		skey := "seckill:goodsStock:10"
		fields := []string{"totalCount", "seckillCount"}
		hmget, err := l.svcCtx.BizRedis.Hmget(skey, fields...)
		if err != nil {
			return nil, err
		}
		totalCount, err := strconv.Atoi(hmget[0])
		if err != nil {
			return nil, err
		}
		seckillCount, err := strconv.Atoi(hmget[1])
		if err != nil {
			return nil, err
		}
		if totalCount-seckillCount > 0 {
			return &types.QuerResponse{
				Message: "正在秒杀",
			}, nil
		} else {
			return &types.QuerResponse{
				Message: "秒杀失败",
			}, nil
		}

	}

}
