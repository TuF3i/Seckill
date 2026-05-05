package handler

import (
	"context"
	"seckill/internal/orderConsumer/core/cache"
	"seckill/internal/orderConsumer/core/dao"
	"sync"

	"gitee.com/liumou_site/logger"
	"github.com/segmentio/kafka-go"
)

type HandlerReliance struct {
	D        *dao.Dao
	C        *cache.Cache
	L        *logger.LocalLogger
	Wg       *sync.WaitGroup
	Ctx      context.Context
	Cancel   context.CancelFunc
	Done     chan struct{}
	Quit     chan struct{}
	Consumer *kafka.Reader
}

type Handler struct {
	*HandlerReliance
}

func NewDao(m *HandlerReliance) *Handler {
	return &Handler{m}
}
