package handler

import (
	"context"
	"errors"
	"seckill/internal/userSvr/core/dto"
	"seckill/internal/userSvr/kitex_gen/usersvr"
	"seckill/pkg/enumTransfer"
	"seckill/pkg/jwt"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func (s *UserSvrImpl) RegisterUser(ctx context.Context, email string, password string) (err error) {
	uid := uuid.New().String()
	if len(email) > 99 || len(email) < 5 {
		return kerrors.NewBizStatusError(dto.InvalidEmail.Status, dto.InvalidEmail.Info)
	}
	if len(password) > 99 || len(password) < 5 {
		return kerrors.NewBizStatusError(dto.InvalidPassword.Status, dto.InvalidPassword.Info)
	}
	err = s.Dao.AddUser(uid, email, password)
	if err != nil {
		return err
	}

	return err
}

func (s *UserSvrImpl) Login(ctx context.Context, email string, password string) (resp *usersvr.JWTToken, err error) {
	if len(email) > 99 || len(email) < 5 {
		return nil, kerrors.NewBizStatusError(dto.InvalidEmail.Status, dto.InvalidEmail.Info)
	}
	if len(password) > 99 || len(password) < 5 {
		return nil, kerrors.NewBizStatusError(dto.InvalidPassword.Status, dto.InvalidPassword.Info)
	}
	data, err := s.Dao.GetUserInfo(email)
	if err != nil {
		return nil, err
	}
	if data.Password != password {
		return nil, kerrors.NewBizStatusError(dto.WrongPassword.Status, dto.WrongPassword.Info)
	}
	accessToken, err := jwt.GenAccessToken(data.Uid, data.Role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwt.GenRefreshToken(data.Uid, data.Role)
	if err != nil {
		return nil, err
	}

	err = s.Cache.AddAccessTokenAndRefreshToken(ctx, data.Uid, accessToken, refreshToken)
	if err != nil {
		return nil, err
	}
	respData := &usersvr.JWTToken{AccessToken: accessToken, RefreshToken: refreshToken}

	return respData, nil
}

func (s *UserSvrImpl) Logout(ctx context.Context, uid string) (err error) {
	_ = s.Cache.DelAccessTokenAndRefreshToken(ctx, uid)
	return nil
}

func (s *UserSvrImpl) RefreshAccessToken(ctx context.Context, claims *usersvr.JWTClaims) (resp string, err error) {
	accessToken, err := jwt.GenAccessToken(claims.UID, claims.Role.String())
	if err != nil {
		return "", err
	}
	err = s.Cache.AddAccessToken(ctx, claims.UID, accessToken)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (s *UserSvrImpl) VerifyAccessToken(ctx context.Context, accessToken string) (resp *usersvr.JWTClaims, err error) {
	claims, err := jwt.VerifyAccessToken(accessToken)
	if err != nil {
		return nil, err
	}
	_accessToken, err := s.Cache.GetAccessToken(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, kerrors.NewBizStatusError(dto.InvalidAccessToken.Status, dto.InvalidAccessToken.Info)
		}
		return nil, err
	}
	if _accessToken != accessToken {
		return nil, kerrors.NewBizStatusError(dto.WrongAccessToken.Status, dto.WrongAccessToken.Info)
	}
	data := &usersvr.JWTClaims{
		UID:  claims.UserID,
		Role: enumTransfer.RoleStringToEnum(claims.Role),
		Type: claims.Type,
	}

	return data, nil
}

func (s *UserSvrImpl) VerifyRefreshToken(ctx context.Context, refreshToken string) (resp *usersvr.JWTClaims, err error) {
	claims, err := jwt.VerifyRefreshToken(refreshToken)
	if err != nil {
		return nil, err
	}
	_refreshToken, err := s.Cache.GetRefreshToken(ctx, claims.UserID)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, kerrors.NewBizStatusError(dto.InvalidRefreshToken.Status, dto.InvalidRefreshToken.Info)
		}
		return nil, err
	}
	if _refreshToken != refreshToken {
		return nil, kerrors.NewBizStatusError(dto.WrongRefreshToken.Status, dto.WrongRefreshToken.Info)
	}
	data := &usersvr.JWTClaims{
		UID:  claims.UserID,
		Role: enumTransfer.RoleStringToEnum(claims.Role),
		Type: claims.Type,
	}

	return data, nil
}
