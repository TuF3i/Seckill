package handler

import (
	"context"
	"errors"
	"seckill/internal/userSvr/core/dto"
	usersvr "seckill/internal/userSvr/kitex_gen/usersvr"
	"seckill/pkg/enumTransfer"
	"seckill/pkg/jwt"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// RegisterUser implements the UserSvrImpl interface.
func (s *UserSvrImpl) RegisterUser(ctx context.Context, email string, password string) (err error) {
	uid := uuid.New().String()
	// 校验数据
	if len(email) > 99 || len(email) < 5 {
		return kerrors.NewBizStatusError(dto.InvalidEmail.Status, dto.InvalidEmail.Info)
	}
	if len(password) > 99 || len(password) < 5 {
		return kerrors.NewBizStatusError(dto.InvalidPassword.Status, dto.InvalidPassword.Info)
	}
	// 调用数据库
	err = s.dao.AddUser(uid, email, password)
	if err != nil {
		return err
	}

	return err
}

// Login implements the UserSvrImpl interface.
func (s *UserSvrImpl) Login(ctx context.Context, email string, password string) (resp *usersvr.JWTToken, err error) {
	// 校验数据
	if len(email) > 99 || len(email) < 5 {
		return nil, kerrors.NewBizStatusError(dto.InvalidEmail.Status, dto.InvalidEmail.Info)
	}
	if len(password) > 99 || len(password) < 5 {
		return nil, kerrors.NewBizStatusError(dto.InvalidPassword.Status, dto.InvalidPassword.Info)
	}
	// 从数据库拉取数据
	data, err := s.dao.GetUserInfo(email)
	if err != nil {
		return nil, err
	}
	// 校验密码
	if data.Password != password {
		return nil, kerrors.NewBizStatusError(dto.WrongPassword.Status, dto.WrongPassword.Info)
	}
	// 生成AccessToken和RefreshToken
	accessToken, err := jwt.GenAccessToken(data.Uid, data.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenRefreshToken(data.Uid, data.Role)
	if err != nil {
		return nil, err
	}

	// 向redis写入Token
	err = s.cache.AddAccessTokenAndRefreshToken(ctx, data.Uid, accessToken, refreshToken)
	if err != nil {
		return nil, err
	}
	// 组装data
	respData := &usersvr.JWTToken{AccessToken: accessToken, RefreshToken: refreshToken}

	return respData, nil
}

// Logout implements the UserSvrImpl interface.
func (s *UserSvrImpl) Logout(ctx context.Context, uid string) (err error) {
	_ = s.cache.DelAccessTokenAndRefreshToken(ctx, uid)
	return nil
}

// RefreshAccessToken implements the UserSvrImpl interface.
func (s *UserSvrImpl) RefreshAccessToken(ctx context.Context, claims *usersvr.JWTClaims) (resp string, err error) {
	// 生成新AccessToken
	accessToken, err := jwt.GenAccessToken(claims.UID, claims.Role.String())
	if err != nil {
		return "", err
	}
	// 写入Redis
	err = s.cache.AddAccessToken(ctx, claims.UID, accessToken)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// VerifyAccessToken implements the UserSvrImpl interface.
func (s *UserSvrImpl) VerifyAccessToken(ctx context.Context, accessToken string) (resp *usersvr.JWTClaims, err error) {
	// 解析AccessToken
	claims, err := jwt.VerifyAccessToken(accessToken)
	if err != nil {
		return nil, err
	}
	// 从cache中校验AccessToken
	_accessToken, err := s.cache.GetAccessToken(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, kerrors.NewBizStatusError(dto.InvalidAccessToken.Status, dto.InvalidAccessToken.Info)
		}
		return nil, err
	}
	if _accessToken != accessToken {
		return nil, kerrors.NewBizStatusError(dto.WrongAccessToken.Status, dto.WrongAccessToken.Info)
	}
	// 组装data
	data := &usersvr.JWTClaims{
		UID:  claims.UserID,
		Role: enumTransfer.RoleStringToEnum(claims.Role),
		Type: claims.Type,
	}

	return data, nil
}

// VerifyRefreshToken implements the UserSvrImpl interface.
func (s *UserSvrImpl) VerifyRefreshToken(ctx context.Context, refreshToken string) (resp *usersvr.JWTClaims, err error) {
	// 解析RefreshToken
	claims, err := jwt.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	// 从cache中校验RefreshToken
	_refreshToken, err := s.cache.GetRefreshToken(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, kerrors.NewBizStatusError(dto.InvalidRefreshToken.Status, dto.InvalidRefreshToken.Info)
		}
		return nil, err
	}
	if _refreshToken != refreshToken {
		return nil, kerrors.NewBizStatusError(dto.WrongRefreshToken.Status, dto.WrongRefreshToken.Info)
	}
	// 组装data
	data := &usersvr.JWTClaims{
		UID:  claims.UserID,
		Role: enumTransfer.RoleStringToEnum(claims.Role),
		Type: claims.Type,
	}

	return data, nil
}
