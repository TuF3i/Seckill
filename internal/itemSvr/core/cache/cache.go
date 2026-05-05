package cache

import (
	"context"
	"seckill/internal/itemSvr/core/pkg/lkeygen"
	"time"

	_ "embed"

	"github.com/redis/go-redis/v9"
)

//go:embed script/order.lua
var prepareOrderScript string

//go:embed script/prepare_order_no_limit.lua
var prepareOrderNoLimitScript string

var prepareOrderCmd = redis.NewScript(prepareOrderScript)
var prepareOrderNoLimitCmd = redis.NewScript(prepareOrderNoLimitScript)

func (r *Cache) PrepareOrderAtomic(ctx context.Context, itemId string, userId string, limitTTL time.Duration) (int64, error) {
	flashKey := lkeygen.GenItemFlashKey(itemId)
	stockKey := lkeygen.GenItemStockKey(itemId)
	limitKey := lkeygen.GenItemPurchaseLimitKey(itemId, userId)

	result, err := prepareOrderCmd.Run(ctx, r.Rdb, []string{flashKey, stockKey, limitKey}, int64(limitTTL.Seconds())).Int64()
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (r *Cache) PrepareOrderAtomicNoLimit(ctx context.Context, itemId string) (int64, error) {
	flashKey := lkeygen.GenItemFlashKey(itemId)
	stockKey := lkeygen.GenItemStockKey(itemId)

	result, err := prepareOrderNoLimitCmd.Run(ctx, r.Rdb, []string{flashKey, stockKey}).Int64()
	if err != nil {
		return 0, err
	}

	return result, nil
}

func (r *Cache) WarmUpItemStock(ctx context.Context, itemId string, stock int64) error {
	key := lkeygen.GenItemStockKey(itemId)
	return r.Rdb.Set(ctx, key, stock, 0).Err()
}

func (r *Cache) GetItemStock(ctx context.Context, itemId string) (int64, error) {
	key := lkeygen.GenItemStockKey(itemId)
	return r.Rdb.Get(ctx, key).Int64()
}

func (r *Cache) DecrItemStock(ctx context.Context, itemId string) (int64, error) {
	key := lkeygen.GenItemStockKey(itemId)
	return r.Rdb.Decr(ctx, key).Result()
}

func (r *Cache) IncrItemStock(ctx context.Context, itemId string) (int64, error) {
	key := lkeygen.GenItemStockKey(itemId)
	return r.Rdb.Incr(ctx, key).Result()
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
