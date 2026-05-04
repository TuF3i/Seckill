package router

import "github.com/cloudwego/hertz/pkg/app/server"

func (r *Router) InitRouter(h *server.Hertz) {
	userGroup := h.Group("/user")
	{
		userGroup.POST("/register", r.HandlerFunc.RegisterUserHandlerFunc())
		userGroup.POST("/login", r.HandlerFunc.LoginHandlerFunc())
		userGroup.GET("/logout", r.Middleware.JWTAuthMiddleware(), r.HandlerFunc.LogoutHandlerFunc())
		userGroup.GET("/refresh", r.Middleware.JWTRefreshMiddleware(), r.HandlerFunc.RefreshAccessTokenHandlerFunc())
	}

	itemGroup := h.Group("/item")
	{
		itemGroup.POST("/add", r.Middleware.JWTAuthMiddleware(), r.HandlerFunc.AddItemHandlerFunc())
		itemGroup.POST("/delete", r.Middleware.JWTAuthMiddleware(), r.HandlerFunc.DeleteItemHandlerFunc())
		itemGroup.GET("/list", r.Middleware.JWTAuthMiddleware(), r.HandlerFunc.ListItemsHandlerFunc())
		itemGroup.POST("/flash/start", r.Middleware.JWTAuthMiddleware(), r.HandlerFunc.StartFlashSaleHandlerFunc())
		itemGroup.POST("/flash/stop", r.Middleware.JWTAuthMiddleware(), r.HandlerFunc.StopFlashSaleHandlerFunc())
	}

	orderGroup := h.Group("/order")
	{
		orderGroup.POST("/create", r.Middleware.JWTAuthMiddleware(), r.HandlerFunc.CreateOrderHandlerFunc())
		orderGroup.GET("/paid", r.Middleware.JWTAuthMiddleware(), r.HandlerFunc.QueryPaidOrdersHandlerFunc())
		orderGroup.GET("/unpaid", r.Middleware.JWTAuthMiddleware(), r.HandlerFunc.QueryUnpaidOrdersHandlerFunc())
		orderGroup.GET("/cancelled", r.Middleware.JWTAuthMiddleware(), r.HandlerFunc.QueryCancelledOrdersHandlerFunc())
	}

	paymentGroup := h.Group("/payment")
	{
		paymentGroup.POST("/process", r.Middleware.JWTAuthMiddleware(), r.HandlerFunc.ProcessPaymentHandlerFunc())
	}
}
