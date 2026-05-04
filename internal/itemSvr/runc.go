package main

import (
	"seckill/infrastructures/postgres"
	"seckill/infrastructures/redis"
	"seckill/internal/itemSvr/core/cache"
	"seckill/internal/itemSvr/core/dao"
	"seckill/internal/itemSvr/core/handler"
	itemSvr "seckill/internal/itemSvr/kitex_gen/itemsvr/itemsvr"
)

var (
	itemSvrObj *handler.ItemSvrImpl
)

const (
	RedisDBItem = 1
)

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

	itemSvrObj = handler.NewItemSvrImpl(&handler.ItemSvrImplReliance{
		Dao:   d,
		Cache: c,
	})
}

func OnDestory() {
	itemSvrObj = nil
}

func main() {
	OnCreate()

	svr := itemSvr.NewServer(itemSvrObj)

	err := svr.Run()
	if err != nil {
		panic(err)
	}

	OnDestory()
}
