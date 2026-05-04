package app

import (
	"context"
	"encoding/json"
	"log"
	"seckill/internal/orderConsumer/core/cache"
	"seckill/internal/orderConsumer/core/dao"
	"seckill/internal/orderConsumer/core/models"
	kafka2 "seckill/infrastructures/kafka"
	"seckill/infrastructures/postgres"
	"seckill/infrastructures/redis"
	"time"
)

type OrderMessage struct {
	OrderId string  `json:"orderId"`
	UserId  string  `json:"userId"`
	ItemId  string  `json:"itemId"`
	Price   float64 `json:"price"`
}

const (
	RedisDBOrder     = 2
	OrderTimeout     = 10 * time.Minute
	TimeoutCheckTick = 1 * time.Minute
)

var (
	d    *dao.Dao
	c    *cache.Cache
	done chan struct{}
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

	d = dao.NewDao(&dao.DaoReliance{
		Pgdb: pgClient,
	})

	c = cache.NewCache(&cache.CacheReliance{
		Rdb: redisClient,
	})

	done = make(chan struct{})
}

func OnDestory() {
	close(done)
	d = nil
	c = nil
}

func RunServer() {
	go timeoutChecker()

	consumer := kafka2.NewKafkaConsumerGroup(
		kafka2.WithTopic("order_topic"),
		kafka2.WithGroupID("order-consumer-group"),
		kafka2.WithBrokers([]string{"localhost:9092"}),
	)

	ctx := context.Background()

	log.Println("OrderConsumer started, waiting for messages...")

	for {
		msg, err := consumer.ReadMessage(ctx)
		if err != nil {
			log.Printf("Error reading message: %v", err)
			continue
		}

		var orderMsg OrderMessage
		err = json.Unmarshal(msg.Value, &orderMsg)
		if err != nil {
			log.Printf("Error unmarshalling message: %v", err)
			continue
		}

		order := &models.Order{
			OrderId: orderMsg.OrderId,
			UserId:  orderMsg.UserId,
			ItemId:  orderMsg.ItemId,
			Price:   orderMsg.Price,
			Status:  1,
		}

		err = d.SaveOrder(order)
		if err != nil {
			log.Printf("Error saving order %s: %v", orderMsg.OrderId, err)
			continue
		}

		err = c.SetOrderStatus(ctx, orderMsg.OrderId, 1)
		if err != nil {
			log.Printf("Error setting order status in cache for %s: %v", orderMsg.OrderId, err)
			continue
		}

		log.Printf("Order %s saved successfully", orderMsg.OrderId)
	}
}

func timeoutChecker() {
	ticker := time.NewTicker(TimeoutCheckTick)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			checkTimeout()
		}
	}
}

func checkTimeout() {
	orders, err := d.GetUnpaidOrders()
	if err != nil {
		log.Printf("Error fetching unpaid orders: %v", err)
		return
	}

	now := time.Now()
	for _, order := range orders {
		if now.Sub(order.CreatedAt) > OrderTimeout {
			err = d.UpdateOrderStatus(order.OrderId, 3)
			if err != nil {
				log.Printf("Error cancelling order %s: %v", order.OrderId, err)
				continue
			}

			err = c.SetOrderStatus(context.Background(), order.OrderId, 3)
			if err != nil {
				log.Printf("Error updating cache for cancelled order %s: %v", order.OrderId, err)
				continue
			}

			log.Printf("Order %s auto-cancelled due to timeout", order.OrderId)
		}
	}
}
