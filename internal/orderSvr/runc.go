package main

import (
	kafka2 "seckill/infrastructures/kafka"
	"seckill/infrastructures/postgres"
	"seckill/infrastructures/redis"
	"seckill/internal/orderSvr/core/cache"
	"seckill/internal/orderSvr/core/dao"
	"seckill/internal/orderSvr/core/handler"
	orderSvr "seckill/internal/orderSvr/kitex_gen/ordersvr/ordersvr"
)

var (
	orderSvrObj *handler.OrderSvrImpl
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

	orderSvrObj = handler.NewOrderSvrImpl(&handler.OrderSvrImplReliance{
		Dao:       d,
		Cache:     c,
		KafkaProd: kafkaProd,
	})
}

func OnDestory() {
	if orderSvrObj != nil && orderSvrObj.KafkaProd != nil {
		orderSvrObj.KafkaProd.Close()
	}
	orderSvrObj = nil
}

func main() {
	OnCreate()

	svr := orderSvr.NewServer(orderSvrObj)

	err := svr.Run()
	if err != nil {
		panic(err)
	}

	OnDestory()
}
