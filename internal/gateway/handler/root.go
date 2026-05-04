package handler

import (
	itemSvr "seckill/internal/itemSvr/kitex_gen/itemsvr/itemsvr"
	orderSvr "seckill/internal/orderSvr/kitex_gen/ordersvr/ordersvr"
	paymentSvr "seckill/internal/paymentSvr/kitex_gen/paymentsvr/paymentsvr"
	userSvr "seckill/internal/userSvr/kitex_gen/usersvr/usersvr"
)

type HandlerReliance struct {
	UserSvr    userSvr.Client
	ItemSvr    itemSvr.Client
	OrderSvr   orderSvr.Client
	PaymentSvr paymentSvr.Client
}

type Handler struct {
	*HandlerReliance
}

func NewHandler(m *HandlerReliance) *Handler {
	return &Handler{m}
}
