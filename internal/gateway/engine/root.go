package engine

import (
	"seckill/configs"
	"seckill/internal/gateway/router"

	"github.com/cloudwego/hertz/pkg/app/server"
)

type RouterReliance struct {
	Router *router.Router
	Config *configs.Config
}

type Engine struct {
	h *server.Hertz
	*RouterReliance
}

func NewEngine(m *RouterReliance) *Engine {
	e := &Engine{RouterReliance: m}
	e.createApiEngine()
	e.Router.InitRouter(e.h)
	return e
}
