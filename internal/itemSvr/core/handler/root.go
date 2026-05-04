package handler

import (
	"seckill/internal/itemSvr/core/cache"
	"seckill/internal/itemSvr/core/dao"
)

type ItemSvrImplReliance struct {
	Dao   *dao.Dao
	Cache *cache.Cache
}

type ItemSvrImpl struct {
	*ItemSvrImplReliance
}

func NewItemSvrImpl(m *ItemSvrImplReliance) *ItemSvrImpl {
	return &ItemSvrImpl{m}
}
