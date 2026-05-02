package handler

import (
	"context"
	"seckill/internal/gateway/dto"
	"seckill/internal/gateway/dto/user"
	"seckill/internal/gateway/pkg/lcontext"
	"seckill/internal/gateway/pkg/lerror"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func (r *Handler) RegisterUserHandlerFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req user.RegisterUserReq
		// 获取POST Body
		if err := c.BindAndValidate(&req); err != nil {
			resp := dto.InternalError(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}
		// 向UserSvr发起调用
		err := r.UserSvr.RegisterUser(ctx, req.Email, req.Password)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		c.JSON(consts.StatusOK, dto.GenFinalResponse(dto.OperationSuccess, nil))
	}
}

func (r *Handler) LoginHandlerFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		var req user.LoginReq
		// 获取POST Body
		if err := c.BindAndValidate(&req); err != nil {
			resp := dto.InternalError(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}
		// 向UserSvr发起调用
		jwtTokens, err := r.UserSvr.Login(ctx, req.Email, req.Password)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		// 转换结构体
		data := user.LoginResp{AccessToken: jwtTokens.AccessToken, RefreshToken: jwtTokens.RefreshToken}

		c.JSON(consts.StatusOK, dto.GenFinalResponse(dto.OperationSuccess, data))
	}
}

func (r *Handler) LogoutHandlerFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从上下文获取claims
		claims, err := lcontext.GetClaimsFromRequestContext(c)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}
		// 向UserSvr发起调用
		err = r.UserSvr.Logout(ctx, claims.UID)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		c.JSON(consts.StatusOK, dto.GenFinalResponse(dto.OperationSuccess, nil))
	}
}

func (r *Handler) RefreshAccessTokenHandlerFunc() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		// 从上下文获取claims
		claims, err := lcontext.GetClaimsFromRequestContext(c)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		// 向UserSvr发起调用
		accessToken, err := r.UserSvr.RefreshAccessToken(ctx, claims)
		if err != nil {
			resp := lerror.GenErrorResponse(err)
			c.JSON(consts.StatusOK, dto.GenFinalResponse(resp, nil))
			return
		}

		// 转换结构体
		data := user.RefreshAccessTokenResp{AccessToken: accessToken}

		c.JSON(consts.StatusOK, dto.GenFinalResponse(dto.OperationSuccess, data))
	}
}
