package engine

import (
	"fmt"

	"github.com/cloudwego/hertz/pkg/app/server"
	prometheus "github.com/hertz-contrib/monitor-prometheus"
)

func (r *Engine) RunApiEngine() {
	// 构造Url
	url := fmt.Sprintf("%v:%v", r.Config.Hertz.ListenAddr, r.Config.Hertz.ListenPort)
	monitorUrl := fmt.Sprintf("%v:%v", r.Config.Hertz.ListenAddr, r.Config.Hertz.MonitoringPort)
	// 创建服务核心
	r.h = server.Default(server.WithHostPorts(url), server.WithTracer(prometheus.NewServerTracer(monitorUrl, "/hertz")))
	
}
