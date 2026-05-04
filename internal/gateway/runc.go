package gateway

import (
	"seckill/internal/gateway/engine"
	"seckill/internal/gateway/handler"
	"seckill/internal/gateway/middleware"
	"seckill/internal/gateway/pkg/lconfig"
	"seckill/internal/gateway/router"
	itemSvr "seckill/internal/itemSvr/kitex_gen/itemsvr/itemsvr"
	orderSvr "seckill/internal/orderSvr/kitex_gen/ordersvr/ordersvr"
	paymentSvr "seckill/internal/paymentSvr/kitex_gen/paymentsvr/paymentsvr"
	userSvr "seckill/internal/userSvr/kitex_gen/usersvr/usersvr"

	"github.com/bwmarrin/snowflake"
)

var (
	GatewayEngine *engine.Engine
)

func OnCreate() {
	snowFlake, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	userClient, err := userSvr.NewClient("UserSvr")
	if err != nil {
		panic(err)
	}

	itemClient, err := itemSvr.NewClient("ItemSvr")
	if err != nil {
		panic(err)
	}

	orderClient, err := orderSvr.NewClient("OrderSvr")
	if err != nil {
		panic(err)
	}

	paymentClient, err := paymentSvr.NewClient("PaymentSvr")
	if err != nil {
		panic(err)
	}

	h := handler.NewHandler(&handler.HandlerReliance{
		UserSvr:    userClient,
		ItemSvr:    itemClient,
		OrderSvr:   orderClient,
		PaymentSvr: paymentClient,
	})

	m := middleware.NewMiddleware(&middleware.MiddlewareReliance{
		UserSvr:   userClient,
		SnowFlake: snowFlake,
	})

	r := router.NewRouter(&router.RouterReliance{
		Middleware:  m,
		HandlerFunc: h,
	})

	cfg := &lconfig.Config{}

	GatewayEngine = engine.NewEngine(&engine.RouterReliance{
		Router: r,
		Config: cfg,
	})

	GatewayEngine.RunApiEngine()

	r.InitRouter(GatewayEngine.Hertz())
}

func OnDestory() {
	GatewayEngine = nil
}

func RunServer() {
	GatewayEngine.Hertz().Spin()
}
