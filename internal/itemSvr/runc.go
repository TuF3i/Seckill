package itemSvr

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"seckill/configs"
	"seckill/infrastructures/nacos"
	"seckill/infrastructures/opentelemetry"
	"seckill/infrastructures/postgres"
	"seckill/infrastructures/redis"
	"seckill/internal/itemSvr/core/cache"
	"seckill/internal/itemSvr/core/dao"
	"seckill/internal/itemSvr/core/handler"
	itemSvr "seckill/internal/itemSvr/kitex_gen/itemsvr/itemsvr"
	"seckill/pkg/config"
	"seckill/pkg/env"

	"gitee.com/liumou_site/logger"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/provider"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"github.com/kitex-contrib/registry-nacos/registry"
)

const RedisDBItem = 1

var (
	itemSvrObj *handler.ItemSvrImpl
	l          *logger.LocalLogger
	p          provider.OtelProvider
)

func RunItemSvr() {
	l = logger.NewLogger(1)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	basicEnv := env.GetEnv()
	onCreate(basicEnv)

	<-quit

	onDestory()
}

func onCreate(env *configs.BasicEnv) {
	l.Modular = "ItemSvr-OnCreate"

	nacosPort, err := strconv.ParseUint(env.NacosPort, 10, 64)
	if err != nil {
		logger.Emer("Convert Port Failed: %v", err.Error())
		os.Exit(1)
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
		os.Exit(1)
	}

	loader, err := config.NewLoader(nacosClient, env.ConfigID, env.ConfigGroup)
	if err != nil {
		logger.Emer("Setup <ConfigLoader> Failed: %v", err.Error())
		os.Exit(1)
	}

	cfg := loader.GetConfig()

	p = opentelemetry.NewProvider(
		opentelemetry.WithEndpoint(cfg.Opentelemetry.ExportEndpoint),
		opentelemetry.WithServiceName("ItemSvr"),
	)

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
		os.Exit(1)
	}

	redisClient, err := redis.NewRedisSentinelClient(
		redis.WithMasterName(cfg.Redis.MasterName),
		redis.WithSentinelAddrs(cfg.Redis.SentinelAddrs),
		redis.WithPassword(cfg.Redis.Password),
		redis.WithSentinelPassword(cfg.Redis.SentinelPassword),
		redis.WithDB(RedisDBItem),
	)
	if err != nil {
		logger.Emer("Setup <Redis> Failed: %v", err.Error())
		os.Exit(1)
	}

	d := dao.NewDao(&dao.DaoReliance{
		Pgdb: pgClient,
	})

	c := cache.NewCache(&cache.CacheReliance{
		Rdb: redisClient,
	})

	itemSvrObj = handler.NewItemSvrImpl(&handler.ItemSvrImplReliance{
		Dao:       d,
		Cache:     c,
		Benchmark: env.Benchmark == "true",
	})

	svr := itemSvr.NewServer(
		itemSvrObj,
		server.WithMetaHandler(transmeta.ServerTTHeaderHandler),
		server.WithRegistry(registry.NewNacosRegistry(nacosClient.NamingClient)),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "ItemSvr"}),
		server.WithSuite(tracing.NewServerSuite()),
	)

	go func() {
		if err := svr.Run(); err != nil {
			logger.Emer("Run <ItemSvr> Failed: %v", err.Error())
			os.Exit(1)
		}
	}()
}

func onDestory() {
	l.Modular = "ItemSvr-OnDestory"
	p.Shutdown(context.Background())
	itemSvrObj = nil
}
