package middleware

import (
	"context"
	"seckill/internal/gateway/pkg/lvar"
	"seckill/pkg/traceid"

	"github.com/cloudwego/hertz/pkg/app"
)

func (r *Middleware) TraceIDMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从请求头里提取TraceID
		traceID := c.Request.Header.Get(lvar.X_TRACE_ID_HEADER)
		// TraceID为空就生成
		if traceID == "" {
			// 生成新TraceID
			traceID := traceid.NewTraceID(r.SnowFlake)
			// 注入TraceID
			ctx = traceid.WriteTraceIDToCtxAndMetainfo(ctx, traceID)
			// 执行后续逻辑
			c.Next(ctx)
			return
		}
		// TraceID不为空就继续
		ctx = traceid.WriteTraceIDToCtxAndMetainfo(ctx, traceID)
		c.Next(ctx)
		return
	}
}
