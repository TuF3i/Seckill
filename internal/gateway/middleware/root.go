package middleware

import (
	"seckill/internal/userSvr/kitex_gen/usersvr/usersvr"

	"github.com/bwmarrin/snowflake"
)

type MiddlewareReliance struct {
	UserSvr   usersvr.Client
	SnowFlake *snowflake.Node
}

type Middleware struct {
	*MiddlewareReliance
}

func NewMiddleware(m *MiddlewareReliance) *Middleware {
	return &Middleware{m}
}
