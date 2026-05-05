package orderSvr

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"seckill/configs"
	"seckill/infrastructures/kafka"
	"seckill/infrastructures/nacos"
	"seckill/infrastructures/postgres"
	"seckill/infrastructures/redis"
	"seckill/internal/orderSvr/core/cache"
	"seckill/internal/orderSvr/core/dao"
	"seckill/internal/orderSvr/core/handler"
	orderSvr "seckill/internal/orderSvr/kitex_gen/ordersvr/ordersvr"
	"seckill/pkg/config"
	"seckill/pkg/env"

	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/registry-nacos/registry"
	"gitee.com/liumou_site/logger"
)

const RedisDBOrder = 2

var (
	orderSvrObj *handler.OrderSvrImpl
	l           *logger.LocalLogger
)

func RunOrderSvr() {
	l = logger.NewLogger(1)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	basicEnv := env.GetEnv()
	onCreate(basicEnv)

	<-quit

	onDestory()
}

func onCreate(env *configs.BasicEnv) {
	l.Modular = "OrderSvr-OnCreate"

	nacosPort, err := strconv.ParseUint(env.NacosPort, 10, 64)
	if err != nil {
		logger.Emer("Convert Port Failed: %v", err.Error())
	}

	nacosClient, err := nacos.NewNacosClient(
		nacos.WithHost(env.NacosAddr),
		nacos.WithPort(nacosPort),
		nacos.WithUserName(env.NacosUser),
		nacos.WithPassword(env.NacosPassword),
		nacos.WithNamespaceID("public"),
	)
	if err != nil {
		logger.Emer("Setup <nacosClient> Failed: %v", err.Error())
	}

	loader, err := config.NewLoader(nacosClient, env.ConfigID, env.ConfigGroup)
	if err != nil {
		logger.Emer("Setup <ConfigLoader> Failed: %v", err.Error())
	}

	cfg := loader.GetConfig()

	portStr := strconv.Itoa(cfg.PostgreSQL.Port)
	pgClient, err := postgres.NewPostgresClient(
		postgres.WithHost(cfg.PostgreSQL.Host),
		postgres.WithPort(portStr),
		postgres.WithUser(cfg.PostgreSQL.User),
		postgres.WithPassword(cfg.PostgreSQL.Password),
		postgres.WithDefaultDB(cfg.PostgreSQL.DefaultDB),
		postgres.WithSSlMode(cfg.PostgreSQL.SSLMode),
	)
	if err != nil {
		logger.Emer("Setup <Postgres> Failed: %v", err.Error())
	}

	redisClient, err := redis.NewRedisSentinelClient(
		redis.WithMasterName(cfg.Redis.MasterName),
		redis.WithSentinelAddrs(cfg.Redis.SentinelAddrs),
		redis.WithPassword(cfg.Redis.Password),
		redis.WithSentinelPassword(cfg.Redis.SentinelPassword),
		redis.WithDB(RedisDBOrder),
	)
	if err != nil {
		logger.Emer("Setup <Redis> Failed: %v", err.Error())
	}

	kafkaProd := kafka.NewKafkaProducerClient(
		kafka.WithTopic("order_topic"),
		kafka.WithBrokers(cfg.Kafka.Brokers),
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

	svr := orderSvr.NewServer(orderSvrObj, server.WithRegistry(registry.NewNacosRegistry(nacosClient.NamingClient)))

	go func() {
		if err := svr.Run(); err != nil {
			logger.Emer("Run <OrderSvr> Failed: %v", err.Error())
		}
	}()
}

func onDestory() {
	l.Modular = "OrderSvr-OnDestory"
	if orderSvrObj != nil && orderSvrObj.KafkaProd != nil {
		orderSvrObj.KafkaProd.Close()
	}
	orderSvrObj = nil
}
