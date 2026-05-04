package app

import (
	"seckill/internal/userSvr/core/cache"
	"seckill/internal/userSvr/core/dao"
	"seckill/internal/userSvr/core/handler"
	userSvr "seckill/internal/userSvr/kitex_gen/usersvr/usersvr"
	"seckill/infrastructures/postgres"
	"seckill/infrastructures/redis"
)

const RedisDBToken = 0

var UserSvrObj *handler.UserSvrImpl

func OnCreate() {
	pgClient, err := postgres.NewPostgresClient()
	if err != nil {
		panic(err)
	}

	redisClient, err := redis.NewRedisSentinelClient(redis.WithDB(RedisDBToken))
	if err != nil {
		panic(err)
	}

	d := dao.NewDao(&dao.DaoReliance{
		Pgdb: pgClient,
	})

	c := cache.NewCache(&cache.CacheReliance{
		Rdb: redisClient,
	})

	UserSvrObj = handler.NewUserSvrImpl(&handler.UserSvrImplReliance{
		Dao:   d,
		Cache: c,
	})
}

func OnDestory() {
	UserSvrObj = nil
}

func RunServer() {
	svr := userSvr.NewServer(UserSvrObj)

	err := svr.Run()
	if err != nil {
		panic(err)
	}
}
