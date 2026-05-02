package handler

import (
	"context"
	usersvr "seckill/internal/userSvr/kitex_gen/usersvr"
)

// UserSvrImpl implements the last service interface defined in the IDL.
type UserSvrImpl struct{}

// RegisterUser implements the UserSvrImpl interface.
func (s *UserSvrImpl) RegisterUser(ctx context.Context, email string, password string) (err error) {
	// TODO: Your code here...
	return
}

// Login implements the UserSvrImpl interface.
func (s *UserSvrImpl) Login(ctx context.Context, email string, password string) (resp *usersvr.JWTToken, err error) {
	// TODO: Your code here...
	return
}

// Logout implements the UserSvrImpl interface.
func (s *UserSvrImpl) Logout(ctx context.Context, uid string) (err error) {
	// TODO: Your code here...
	return
}

// RefreshAccessToken implements the UserSvrImpl interface.
func (s *UserSvrImpl) RefreshAccessToken(ctx context.Context, claims *usersvr.JWTClaims) (resp string, err error) {
	// TODO: Your code here...
	return
}

// VerifyAccessToken implements the UserSvrImpl interface.
func (s *UserSvrImpl) VerifyAccessToken(ctx context.Context, accessToken string) (resp *usersvr.JWTClaims, err error) {
	// TODO: Your code here...
	return
}

// VerifyRefreshToken implements the UserSvrImpl interface.
func (s *UserSvrImpl) VerifyRefreshToken(ctx context.Context, refreshToken string) (resp *usersvr.JWTClaims, err error) {
	// TODO: Your code here...
	return
}
