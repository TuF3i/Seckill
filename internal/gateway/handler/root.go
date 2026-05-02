package handler

import "seckill/internal/userSvr/kitex_gen/usersvr/usersvr"

type HandlerReliance struct {
	UserSvr usersvr.Client
}

type Handler struct {
	*HandlerReliance
}

func NewHandler(m *HandlerReliance) *Handler {
	return &Handler{m}
}
