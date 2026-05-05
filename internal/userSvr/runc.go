package userSvr

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"seckill/configs"
	"seckill/infrastructures/nacos"
	"seckill/infrastructures/postgres"
	"seckill/infrastructures/redis"
	"seckill/internal/userSvr/core/cache"
	"seckill/internal/userSvr/core/dao"
	"seckill/internal/userSvr/core/handler"
	userSvr "seckill/internal/userSvr/kitex_gen/usersvr/usersvr"
	"seckill/pkg/config"
	"seckill/pkg/env"

	"gitee.com/liumou_site/logger"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/registry-nacos/registry"
)

const RedisDBToken = 0

var (
	userSvrObj *handler.UserSvrImpl
	l          *logger.LocalLogger
)

func RunUserSvr() {
	l = logger.NewLogger(1)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	basicEnv := env.GetEnv()
	onCreate(basicEnv)

	<-quit

	onDestory()
}

func onCreate(env *configs.BasicEnv) {
	l.Modular = "UserSvr-OnCreate"

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
		redis.WithDB(RedisDBToken),
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

	userSvrObj = handler.NewUserSvrImpl(&handler.UserSvrImplReliance{
		Dao:   d,
		Cache: c,
	})

	svr := userSvr.NewServer(
		userSvrObj,
		server.WithRegistry(registry.NewNacosRegistry(nacosClient.NamingClient)),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "UserSvr"}),
	)

	go func() {
		if err := svr.Run(); err != nil {
			logger.Emer("Run <UserSvr> Failed: %v", err.Error())
			os.Exit(1)
		}
	}()
}

func onDestory() {
	l.Modular = "UserSvr-OnDestory"
	userSvrObj = nil
}
