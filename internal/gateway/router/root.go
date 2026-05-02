package router

import (
	"seckill/internal/gateway/handler"
	"seckill/internal/gateway/middleware"
)

type RouterReliance struct {
	Middleware  *middleware.Middleware
	HandlerFunc *handler.Handler
}

type Router struct {
	*RouterReliance
}

func NewRouter(m *RouterReliance) *Router {
	return &Router{m}
}
