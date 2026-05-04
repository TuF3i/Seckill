package app

import (
	"seckill/internal/itemSvr/core/cache"
	"seckill/internal/itemSvr/core/dao"
	"seckill/internal/itemSvr/core/handler"
	itemSvr "seckill/internal/itemSvr/kitex_gen/itemsvr/itemsvr"
	"seckill/infrastructures/postgres"
	"seckill/infrastructures/redis"
)

const RedisDBItem = 1

var ItemSvrObj *handler.ItemSvrImpl

func OnCreate() {
	pgClient, err := postgres.NewPostgresClient()
	if err != nil {
		panic(err)
	}

	redisClient, err := redis.NewRedisSentinelClient(redis.WithDB(RedisDBItem))
	if err != nil {
		panic(err)
	}

	d := dao.NewDao(&dao.DaoReliance{
		Pgdb: pgClient,
	})

	c := cache.NewCache(&cache.CacheReliance{
		Rdb: redisClient,
	})

	ItemSvrObj = handler.NewItemSvrImpl(&handler.ItemSvrImplReliance{
		Dao:   d,
		Cache: c,
	})
}

func OnDestory() {
	ItemSvrObj = nil
}

func RunServer() {
	svr := itemSvr.NewServer(ItemSvrObj)

	err := svr.Run()
	if err != nil {
		panic(err)
	}
}
