package cache

import (
	"context"
	"seckill/internal/userSvr/core/pkg/lkeygen"
	"seckill/pkg/jwt"
)

func (r *Cache) AddAccessTokenAndRefreshToken(ctx context.Context, uid string, accessToken string, refreshToken string) error {
	pipe := r.Rdb.TxPipeline()

	key := lkeygen.GenAccessTokenKey(uid)
	pipe.Set(ctx, key, accessToken, jwt.GetAccessTokenExpireTime())

	key = lkeygen.GenRefreshTokenKey(uid)
	pipe.Set(ctx, key, refreshToken, jwt.GetRefreshTokenExpireTime())

	_, err := pipe.Exec(ctx)

	return err
}

func (r *Cache) AddAccessToken(ctx context.Context, uid string, token string) error {
	key := lkeygen.GenAccessTokenKey(uid)

	return r.Rdb.Set(ctx, key, token, jwt.GetAccessTokenExpireTime()).Err()
}

func (r *Cache) AddRefreshToken(ctx context.Context, uid string, token string) error {
	key := lkeygen.GenRefreshTokenKey(uid)

	return r.Rdb.Set(ctx, key, token, jwt.GetRefreshTokenExpireTime()).Err()
}

func (r *Cache) GetAccessToken(ctx context.Context, uid string) (string, error) {
	key := lkeygen.GenAccessTokenKey(uid)
	return r.Rdb.Get(ctx, key).Result()
}

func (r *Cache) GetRefreshToken(ctx context.Context, uid string) (string, error) {
	key := lkeygen.GenRefreshTokenKey(uid)
	return r.Rdb.Get(ctx, key).Result()
}

func (r *Cache) ExistsAccessToken(ctx context.Context, uid string) (bool, error) {
	key := lkeygen.GenAccessTokenKey(uid)
	n, err := r.Rdb.Exists(ctx, key).Result()
	return n > 0, err
}

func (r *Cache) ExistsRefreshToken(ctx context.Context, uid string) (bool, error) {
	key := lkeygen.GenRefreshTokenKey(uid)
	n, err := r.Rdb.Exists(ctx, key).Result()
	return n > 0, err
}

func (r *Cache) DelAccessTokenAndRefreshToken(ctx context.Context, uid string) error {
	pipe := r.Rdb.TxPipeline()

	key := lkeygen.GenAccessTokenKey(uid)
	pipe.Del(ctx, key)

	key = lkeygen.GenRefreshTokenKey(uid)
	pipe.Del(ctx, key)

	_, err := pipe.Exec(ctx)

	return err
}
