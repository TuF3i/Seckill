package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app/server"
	prometheus "github.com/hertz-contrib/monitor-prometheus"
)

func (r *Engine) RunApiEngine() {
	// 构造Url
	url := fmt.Sprintf("%v:%v", r.Config.Gateway.ListenAddr, r.Config.Gateway.ListenPort)
	monitorUrl := fmt.Sprintf("%v:%v", r.Config.Gateway.ListenAddr, r.Config.Gateway.MonitoringPort)
	// 创建服务核心
	r.h = server.Default(server.WithHostPorts(url), server.WithTracer(prometheus.NewServerTracer(monitorUrl, "/hertz")))
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
