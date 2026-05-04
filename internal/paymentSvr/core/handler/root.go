package handler

import (
	"seckill/internal/paymentSvr/core/cache"
	"seckill/internal/paymentSvr/core/dao"
)

type PaymentSvrImplReliance struct {
	Dao   *dao.Dao
	Cache *cache.Cache
}

type PaymentSvrImpl struct {
	*PaymentSvrImplReliance
}

func NewPaymentSvrImpl(m *PaymentSvrImplReliance) *PaymentSvrImpl {
	return &PaymentSvrImpl{m}
}
