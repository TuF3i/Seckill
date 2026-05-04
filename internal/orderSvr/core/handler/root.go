package handler

import (
	"seckill/internal/orderSvr/core/cache"
	"seckill/internal/orderSvr/core/dao"
	"github.com/segmentio/kafka-go"
)

type OrderSvrImplReliance struct {
	Dao        *dao.Dao
	Cache      *cache.Cache
	KafkaProd  *kafka.Writer
}

type OrderSvrImpl struct {
	*OrderSvrImplReliance
}

func NewOrderSvrImpl(m *OrderSvrImplReliance) *OrderSvrImpl {
	return &OrderSvrImpl{m}
}
