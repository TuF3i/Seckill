package app

import (
	"seckill/internal/paymentSvr/core/cache"
	"seckill/internal/paymentSvr/core/dao"
	"seckill/internal/paymentSvr/core/handler"
	paymentSvr "seckill/internal/paymentSvr/kitex_gen/paymentsvr/paymentsvr"
	"seckill/infrastructures/postgres"
	"seckill/infrastructures/redis"
)

const RedisDBOrder = 2

var PaymentSvrObj *handler.PaymentSvrImpl

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

	PaymentSvrObj = handler.NewPaymentSvrImpl(&handler.PaymentSvrImplReliance{
		Dao:   d,
		Cache: c,
	})
}

func OnDestory() {
	PaymentSvrObj = nil
}

func RunServer() {
	svr := paymentSvr.NewServer(PaymentSvrObj)

	err := svr.Run()
	if err != nil {
		panic(err)
	}
}
