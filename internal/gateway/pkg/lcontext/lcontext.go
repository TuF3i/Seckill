package lcontext

import (
	"errors"
	"seckill/internal/userSvr/kitex_gen/usersvr"

	"github.com/cloudwego/hertz/pkg/app"
)

const (
	JWT_CLAIMS_CTX_KEY = "jwt.claims"
	JWT_TOKEN_CTX_KEY  = "jwt.token"
)

func PutClaimsInRequestContext(c *app.RequestContext, claims *usersvr.JWTClaims) {
	c.Set(JWT_CLAIMS_CTX_KEY, claims)
}

func GetClaimsFromRequestContext(c *app.RequestContext) (*usersvr.JWTClaims, error) {
	raw, ok := c.Get(JWT_CLAIMS_CTX_KEY)
	if !ok {
		return nil, errors.New("claims not found")
	}

	claims, ok := raw.(*usersvr.JWTClaims)
	if !ok {
		return nil, errors.New("type assertion failed")
	}

	return claims, nil
}

func PutTokenInRequestContext(c *app.RequestContext, token string) {
	c.Set(JWT_TOKEN_CTX_KEY, token)
}

func GetTokenFromRequestContext(c *app.RequestContext) (string, error) {
	raw, ok := c.Get(JWT_TOKEN_CTX_KEY)
	if !ok {
		return "", errors.New("token not found")
	}

	token, ok := raw.(string)
	if !ok {
		return "", errors.New("type assertion failed")
	}

	return token, nil
}
