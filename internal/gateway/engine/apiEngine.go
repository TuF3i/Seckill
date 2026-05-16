package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	prometheus "github.com/hertz-contrib/monitor-prometheus"
	"github.com/hertz-contrib/obs-opentelemetry/tracing"
)

func (r *Engine) createApiEngine() {
	// 构造Url
	url := fmt.Sprintf("%v:%v", r.Config.Gateway.ListenAddr, r.Config.Gateway.ListenPort)
	monitorUrl := fmt.Sprintf("%v:%v", r.Config.Gateway.ListenAddr, r.Config.Gateway.MonitoringPort)
	// 创建链路追踪
	tracer, cfg := tracing.NewServerTracer()
	// 创建服务核心
	r.h = server.Default(
		tracer,
		server.WithHostPorts(url),
		server.WithTracer(prometheus.NewServerTracer(monitorUrl, "/hertz")),
	)
	// 使用全局追踪中间件
	r.h.Use(tracing.ServerMiddleware(cfg))
}

func (r *Engine) RunApiEngine() {
	// 启动服务核心
	go func() { r.h.Spin() }()
}

func (r *Engine) StopApiEngine() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := r.h.Shutdown(ctx); err != nil { // 会触发优雅停服
		return err
	}
	return nil
}
