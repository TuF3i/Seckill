package main

import (
	"seckill/infrastructures/postgres"
	"seckill/infrastructures/redis"
	"seckill/internal/paymentSvr/core/cache"
	"seckill/internal/paymentSvr/core/dao"
	"seckill/internal/paymentSvr/core/handler"
	paymentSvr "seckill/internal/paymentSvr/kitex_gen/paymentsvr/paymentsvr"
)

var (
	paymentSvrObj *handler.PaymentSvrImpl
)

const (
	RedisDBOrder = 2
)

func OnCreate() {
	pgClient, err := postgres.NewPostgresClient()
	if err != nil {
		panic(err)
	}

	redisClient, err := redis.NewRedisSentinelClient(redis.WithDB(RedisDBOrder))
	if err != nil {
		panic(err)
	}

	d := dao.NewDao(&dao.DaoReliance{
		Pgdb: pgClient,
	})

	c := cache.NewCache(&cache.CacheReliance{
		Rdb: redisClient,
	})

	paymentSvrObj = handler.NewPaymentSvrImpl(&handler.PaymentSvrImplReliance{
		Dao:   d,
		Cache: c,
	})
}

func OnDestory() {
	paymentSvrObj = nil
}

func main() {
	OnCreate()

	svr := paymentSvr.NewServer(paymentSvrObj)

	err := svr.Run()
	if err != nil {
		panic(err)
	}

	OnDestory()
}
