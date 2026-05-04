package cache

import (
	"context"
	"seckill/internal/itemSvr/core/pkg/lkeygen"
	"time"
)

func (r *Cache) WarmUpItemStock(ctx context.Context, itemId string, stock int64) error {
	key := lkeygen.GenItemStockKey(itemId)
	return r.Rdb.Set(ctx, key, stock, 0).Err()
}

func (r *Cache) GetItemStock(ctx context.Context, itemId string) (int64, error) {
	key := lkeygen.GenItemStockKey(itemId)
	return r.Rdb.Get(ctx, key).Int64()
}

func (r *Cache) SetFlashStatus(ctx context.Context, itemId string, status int32) error {
	key := lkeygen.GenItemFlashKey(itemId)
	return r.Rdb.Set(ctx, key, status, 0).Err()
}

func (r *Cache) GetFlashStatus(ctx context.Context, itemId string) (int32, error) {
	key := lkeygen.GenItemFlashKey(itemId)
	val, err := r.Rdb.Get(ctx, key).Int64()
	if err != nil {
		return 0, err
	}
	return int32(val), nil
}

func (r *Cache) SetPurchaseLimit(ctx context.Context, userId string, itemId string) error {
	key := lkeygen.GenItemPurchaseLimitKey(itemId, userId)
	return r.Rdb.Set(ctx, key, 1, 5*time.Minute).Err()
}

func (r *Cache) ExistsPurchaseLimit(ctx context.Context, userId string, itemId string) (bool, error) {
	key := lkeygen.GenItemPurchaseLimitKey(itemId, userId)
	n, err := r.Rdb.Exists(ctx, key).Result()
	return n > 0, err
}

func (r *Cache) DelItemFlashCache(ctx context.Context, itemId string) error {
	pipe := r.Rdb.TxPipeline()

	stockKey := lkeygen.GenItemStockKey(itemId)
	pipe.Del(ctx, stockKey)

	flashKey := lkeygen.GenItemFlashKey(itemId)
	pipe.Del(ctx, flashKey)

	_, err := pipe.Exec(ctx)
	return err
}
