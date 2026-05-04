package app

import (
	kafka2 "seckill/infrastructures/kafka"
	"seckill/infrastructures/postgres"
	"seckill/infrastructures/redis"
	"seckill/internal/orderSvr/core/cache"
	"seckill/internal/orderSvr/core/dao"
	"seckill/internal/orderSvr/core/handler"
	orderSvr "seckill/internal/orderSvr/kitex_gen/ordersvr/ordersvr"
)

const RedisDBOrder = 2

var OrderSvrObj *handler.OrderSvrImpl

func OnCreate() {
	pgClient, err := postgres.NewPostgresClient()
	if err != nil {
		panic(err)
	}

	redisClient, err := redis.NewRedisSentinelClient(redis.WithDB(RedisDBOrder))
	if err != nil {
		panic(err)
	}

	kafkaProd := kafka2.NewKafkaProducerClient(
		kafka2.WithTopic("order_topic"),
		kafka2.WithBrokers([]string{"localhost:9092"}),
	)

	d := dao.NewDao(&dao.DaoReliance{
		Pgdb: pgClient,
	})

	c := cache.NewCache(&cache.CacheReliance{
		Rdb: redisClient,
	})

	OrderSvrObj = handler.NewOrderSvrImpl(&handler.OrderSvrImplReliance{
		Dao:       d,
		Cache:     c,
		KafkaProd: kafkaProd,
	})
}

func OnDestory() {
	if OrderSvrObj != nil && OrderSvrObj.KafkaProd != nil {
		OrderSvrObj.KafkaProd.Close()
	}
	OrderSvrObj = nil
}

func RunServer() {
	svr := orderSvr.NewServer(OrderSvrObj)

	err := svr.Run()
	if err != nil {
		panic(err)
	}
}
