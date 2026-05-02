package middleware

import "seckill/internal/userSvr/kitex_gen/usersvr/usersvr"

type MiddlewareReliance struct {
	UserSvr usersvr.Client
}

type Middleware struct {
	*MiddlewareReliance
}

func NewMiddleware(m *MiddlewareReliance) *Middleware {
	return &Middleware{m}
}
