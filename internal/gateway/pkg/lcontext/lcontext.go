package lcontext

import (
	"errors"
	"seckill/internal/userSvr/kitex_gen/usersvr"

	"github.com/cloudwego/hertz/pkg/app"
)

const (
	JWT_CLAIMS_CTX_KEY = "jwt.claims"
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
