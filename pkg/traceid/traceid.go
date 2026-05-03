package traceid

import (
	"context"
	"errors"

	"github.com/bwmarrin/snowflake"
	"github.com/bytedance/gopkg/cloud/metainfo"
)

const (
	TRACE_ID_CTX_KEY       = "key.ctx.traceid"
	TRACE_ID_META_INFO_KEY = "key.metainfo.traceid"
)

func NewTraceID(s *snowflake.Node) string {
	return s.Generate().String()
}

func WriteTraceIDToContext(ctx context.Context, traceid string) context.Context {
	return context.WithValue(ctx, TRACE_ID_CTX_KEY, traceid)
}

func GetTraceIDFromContext(ctx context.Context) (string, error) {
	data, ok := ctx.Value(TRACE_ID_CTX_KEY).(string)
	if !ok {
		return "", errors.New("type assertion failed")
	}

	return data, nil
}

func WriteTraceIDToMetainfo(ctx context.Context, traceid string) context.Context {
	return metainfo.WithPersistentValue(ctx, TRACE_ID_META_INFO_KEY, traceid)
}

func GetTraceIDFromMetainfo(ctx context.Context) (string, error) {
	traceID, ok := metainfo.GetPersistentValue(ctx, TRACE_ID_META_INFO_KEY)
	if !ok {
		return "", errors.New("traceid not found")
	}

	return traceID, nil
}

func WriteTraceIDToCtxAndMetainfo(ctx context.Context, traceid string) context.Context {
	ctx = context.WithValue(ctx, TRACE_ID_CTX_KEY, traceid)
	ctx = metainfo.WithPersistentValue(ctx, TRACE_ID_META_INFO_KEY, traceid)
	return ctx
}
