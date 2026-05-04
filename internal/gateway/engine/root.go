package engine

import (
	"seckill/internal/gateway/pkg/lconfig"
	"seckill/internal/gateway/router"

	"github.com/cloudwego/hertz/pkg/app/server"
)

type RouterReliance struct {
	Router *router.Router
	Config *lconfig.Config
}

type Engine struct {
	h *server.Hertz
	*RouterReliance
}

func NewEngine(m *RouterReliance) *Engine {
	return &Engine{RouterReliance: m}
}

func (r *Engine) Hertz() *server.Hertz {
	return r.h
}
