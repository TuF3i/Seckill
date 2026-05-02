package router

import "github.com/cloudwego/hertz/pkg/app/server"

func (r *Router) InitRouter(h *server.Hertz) {
	// 用户路由组
	g := h.Group("/user")
	{
		// 注册
		g.POST("/register", r.HandlerFunc.RegisterUserHandlerFunc())
		// 登录
		g.POST("/login", r.HandlerFunc.LoginHandlerFunc())
		// 退出登录
		g.GET("/logout", r.Middleware.JWTAuthMiddleware(), r.HandlerFunc.LogoutHandlerFunc())
		// 刷新AccessToken
		g.GET("/refresh", r.Middleware.JWTRefreshMiddleware(), r.HandlerFunc.RefreshAccessTokenHandlerFunc())
	}
}
