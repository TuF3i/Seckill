package handler

import (
	"seckill/internal/userSvr/core/cache"
	"seckill/internal/userSvr/core/dao"
)

type UserSvrImplReliance struct {
	Dao   *dao.Dao
	Cache *cache.Cache
}

type UserSvrImpl struct {
	*UserSvrImplReliance
}

func NewUserSvrImpl(m *UserSvrImplReliance) *UserSvrImpl {
	return &UserSvrImpl{m}
}
