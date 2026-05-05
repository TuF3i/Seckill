package orderConsumer

import (
	"context"
	"os"
	"os/signal"
	"seckill/internal/orderConsumer/core/handler"
	"strconv"
	"sync"
	"syscall"

	"seckill/configs"
	"seckill/infrastructures/kafka"
	"seckill/infrastructures/nacos"
	"seckill/infrastructures/postgres"
	"seckill/infrastructures/redis"
	"seckill/internal/orderConsumer/core/cache"
	"seckill/internal/orderConsumer/core/dao"
	"seckill/pkg/config"
	"seckill/pkg/env"

	"gitee.com/liumou_site/logger"
)

const (
	RedisDBOrder = 2
)

var (
	l          *logger.LocalLogger
	handlerObj *handler.Handler
)

func RunOrderConsumer() {
	l = logger.NewLogger(1)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	basicEnv := env.GetEnv()
	onCreate(basicEnv)

	<-quit

	onDestory()
}

func onCreate(env *configs.BasicEnv) {
	l.Modular = "OrderConsumer-OnCreate"

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

	d := dao.NewDao(&dao.DaoReliance{
		Pgdb: pgClient,
	})

	c := cache.NewCache(&cache.CacheReliance{
		Rdb: redisClient,
	})

	consumer := kafka.NewKafkaConsumerGroup(
		kafka.WithTopic("order_topic"),
		kafka.WithGroupID("order-consumer-group"),
		kafka.WithBrokers(cfg.Kafka.Brokers),
	)

	ctx, cancel := context.WithCancel(context.Background())

	handlerObj = &handler.Handler{&handler.HandlerReliance{
		D:        d,
		C:        c,
		L:        logger.NewLogger(1),
		Wg:       &sync.WaitGroup{},
		Ctx:      ctx,
		Cancel:   cancel,
		Done:     make(chan struct{}),
		Quit:     make(chan struct{}),
		Consumer: consumer,
	}}

	handlerObj.RunTimeoutChecker()
	handlerObj.RunOrderConsumer()

	logger.Info("OrderConsumer started, waiting for messages...")

}

func onDestory() {
	l.Modular = "OrderConsumer-OnDestory"
	handlerObj.StopOrderConsumer()
	handlerObj.StopTimeoutChecker()

	handlerObj.WaitUntilDone()
}
