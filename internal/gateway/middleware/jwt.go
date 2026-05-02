package middleware

import (
	"context"
	"seckill/internal/gateway/dto"
	"seckill/internal/gateway/pkg/lcontext"
	"seckill/internal/gateway/pkg/lerror"
	"seckill/pkg/jwt"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func (r *Middleware) JWTAuthMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从请求头提取Authorization-Header
		authHeader := string(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.JSON(consts.StatusOK, dto.EmptyJWTString)
			c.Abort()
			return
		}
		// 提取Token
		token := jwt.StripBearer(authHeader)
		// 向UserSvr校验Token
		claims, err := r.UserSvr.VerifyAccessToken(ctx, token)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			c.Abort()
			return
		}
		// 将claims写入上下文
		lcontext.PutClaimsInRequestContext(c, claims)
	}
}

func (r *Middleware) JWTRefreshMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从请求头提取Authorization-Header
		authHeader := string(c.GetHeader("Authorization"))
		if authHeader == "" {
			c.JSON(consts.StatusOK, dto.EmptyJWTString)
			c.Abort()
			return
		}
		// 提取Token
		token := jwt.StripBearer(authHeader)
		// 向UserSvr校验Token
		claims, err := r.UserSvr.VerifyRefreshToken(ctx, token)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			c.Abort()
			return
		}
		// 将claims写入上下文
		lcontext.PutClaimsInRequestContext(c, claims)
	}
}
