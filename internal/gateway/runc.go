package gateway

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"seckill/configs"
	"seckill/infrastructures/nacos"
	"seckill/internal/gateway/engine"
	"seckill/internal/gateway/handler"
	"seckill/internal/gateway/middleware"
	"seckill/internal/gateway/router"
	itemSvr "seckill/internal/itemSvr/kitex_gen/itemsvr/itemsvr"
	orderSvr "seckill/internal/orderSvr/kitex_gen/ordersvr/ordersvr"
	paymentSvr "seckill/internal/paymentSvr/kitex_gen/paymentsvr/paymentsvr"
	userSvr "seckill/internal/userSvr/kitex_gen/usersvr/usersvr"
	"seckill/pkg/config"
	"seckill/pkg/env"
	"seckill/pkg/stringToNodeID"

	"github.com/kitex-contrib/registry-nacos/resolver"

	rpcclient "github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"

	"gitee.com/liumou_site/logger"
	"github.com/bwmarrin/snowflake"
)

var (
	GatewayEngine *engine.Engine
	l             *logger.LocalLogger
)

func RunGateway() {
	// 创建日志实例
	l = logger.NewLogger(1)

	// 创建信号监听器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 获取环境变量
	basicEnv := env.GetEnv()

	onCreate(basicEnv)

	<-quit

	onDestory()
}

func onCreate(env *configs.BasicEnv) {
	l.Modular = "Gateway-OnCreate"

	// 生成SnowFlake
	snowFlake, err := snowflake.NewNode(stringToNodeID.StringToNodeID(env.ContainerName))
	if err != nil {
		logger.Emer("Setup <snowflake> Failed: %v", err.Error())
		os.Exit(1)
	}

	// 转换数据类型
	nacosPort, err := strconv.ParseUint(env.NacosPort, 10, 64)
	if err != nil {
		logger.Emer("Convert Port Failed: %v", err.Error())
		os.Exit(1)
	}

	// 初始化Nacos
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

	// 初始化userSvr
	userClient, err := userSvr.NewClient(
		"UserSvr",
		rpcclient.WithResolver(resolver.NewNacosResolver(nacosClient.NamingClient)),
		rpcclient.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		rpcclient.WithTransportProtocol(transport.TTHeader),
	)
	if err != nil {
		logger.Emer("Setup <userSvr> Failed: %v", err.Error())
		os.Exit(1)
	}

	// 初始化itemSvr
	itemClient, err := itemSvr.NewClient(
		"ItemSvr",
		rpcclient.WithResolver(resolver.NewNacosResolver(nacosClient.NamingClient)),
		rpcclient.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		rpcclient.WithTransportProtocol(transport.TTHeader),
	)
	if err != nil {
		logger.Emer("Setup <itemSvr> Failed: %v", err.Error())
		os.Exit(1)
	}

	// 初始化orderSvr
	orderClient, err := orderSvr.NewClient(
		"OrderSvr",
		rpcclient.WithResolver(resolver.NewNacosResolver(nacosClient.NamingClient)),
		rpcclient.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		rpcclient.WithTransportProtocol(transport.TTHeader),
	)
	if err != nil {
		logger.Emer("Setup <orderSvr> Failed: %v", err.Error())
		os.Exit(1)
	}

	// 初始化paymentSvr
	paymentClient, err := paymentSvr.NewClient(
		"PaymentSvr",
		rpcclient.WithResolver(resolver.NewNacosResolver(nacosClient.NamingClient)),
		rpcclient.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		rpcclient.WithTransportProtocol(transport.TTHeader),
	)
	if err != nil {
		logger.Emer("Setup <paymentSvr> Failed: %v", err.Error())
		os.Exit(1)
	}

	// 初始化Handler
	h := handler.NewHandler(&handler.HandlerReliance{
		UserSvr:    userClient,
		ItemSvr:    itemClient,
		OrderSvr:   orderClient,
		PaymentSvr: paymentClient,
	})

	// 初始化中间件
	m := middleware.NewMiddleware(&middleware.MiddlewareReliance{
		UserSvr:   userClient,
		SnowFlake: snowFlake,
	})

	// 初始化路由
	r := router.NewRouter(&router.RouterReliance{
		Middleware:  m,
		HandlerFunc: h,
	})

	// 初始化配置信息
	loader, err := config.NewLoader(nacosClient, env.ConfigID, env.ConfigGroup)
	if err != nil {
		logger.Emer("Setup <ConfigLoader> Failed: %v", err.Error())
		os.Exit(1)
	}

	// 初始化API引擎
	GatewayEngine = engine.NewEngine(&engine.RouterReliance{
		Router: r,
		Config: loader.GetConfig(),
	})

	// 运行引擎
	GatewayEngine.RunApiEngine()
}

func onDestory() {
	l.Modular = "Gateway-OnDestory"

	err := GatewayEngine.StopApiEngine()
	if err != nil {
		logger.Warn("Stop <GatewayEngine> Error: %v", err.Error())
	}

	return
}
