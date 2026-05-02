package middleware

import (
	"context"

	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/hertz/pkg/app"
)

func TraceIDMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从请求头里提取TraceID
		traceID := c.Request.Header.Get(union_var.X_TRACE_ID_HEADER)
		// TraceID为空就生成
		if traceID == "" {
			traceID := core.SnowFlake.Generate().String()
			// 向MataInfo中写入TraceID
			ctx = metainfo.WithPersistentValue(ctx, union_var.TRACE_ID_KEY, traceID)
			// 向Context中写入TraceID
			ctx = context.WithValue(ctx, union_var.TRACE_ID_KEY, traceID)
			// 执行后续逻辑
			c.Next(ctx)
			return
		}
		// 向MataInfo中写入TraceID
		ctx = metainfo.WithPersistentValue(ctx, union_var.TRACE_ID_KEY, traceID)
		// 向Context中写入TraceID
		ctx = context.WithValue(ctx, union_var.TRACE_ID_KEY, traceID)
		// 执行后续逻辑
		c.Next(ctx)
	}
}
