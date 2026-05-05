package paymentSvr

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"seckill/configs"
	"seckill/infrastructures/nacos"
	"seckill/infrastructures/postgres"
	"seckill/infrastructures/redis"
	"seckill/internal/paymentSvr/core/cache"
	"seckill/internal/paymentSvr/core/dao"
	"seckill/internal/paymentSvr/core/handler"
	paymentSvr "seckill/internal/paymentSvr/kitex_gen/paymentsvr/paymentsvr"
	"seckill/pkg/config"
	"seckill/pkg/env"

	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/registry-nacos/registry"
	"gitee.com/liumou_site/logger"
)

const RedisDBOrder = 2

var (
	paymentSvrObj *handler.PaymentSvrImpl
	l             *logger.LocalLogger
)

func RunPaymentSvr() {
	l = logger.NewLogger(1)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	basicEnv := env.GetEnv()
	onCreate(basicEnv)

	<-quit

	onDestory()
}

func onCreate(env *configs.BasicEnv) {
	l.Modular = "PaymentSvr-OnCreate"

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

	paymentSvrObj = handler.NewPaymentSvrImpl(&handler.PaymentSvrImplReliance{
		Dao:   d,
		Cache: c,
	})

	svr := paymentSvr.NewServer(paymentSvrObj, server.WithRegistry(registry.NewNacosRegistry(nacosClient.NamingClient)))

	go func() {
		if err := svr.Run(); err != nil {
			logger.Emer("Run <PaymentSvr> Failed: %v", err.Error())
		}
	}()
}

func onDestory() {
	l.Modular = "PaymentSvr-OnDestory"
	paymentSvrObj = nil
}
