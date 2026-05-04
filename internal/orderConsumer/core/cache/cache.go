package cache

import (
	"context"
	"seckill/internal/orderConsumer/core/pkg/lkeygen"
)

func (r *Cache) SetOrderStatus(ctx context.Context, orderId string, status int32) error {
	key := lkeygen.GenOrderStatusKey(orderId)
	return r.Rdb.Set(ctx, key, status, 0).Err()
}

func (r *Cache) SetOrderInfo(ctx context.Context, orderId string, info string) error {
	key := lkeygen.GenOrderInfoKey(orderId)
	return r.Rdb.Set(ctx, key, info, 0).Err()
}

func (r *Cache) DelOrderCache(ctx context.Context, orderId string) error {
	pipe := r.Rdb.TxPipeline()

	statusKey := lkeygen.GenOrderStatusKey(orderId)
	pipe.Del(ctx, statusKey)

	infoKey := lkeygen.GenOrderInfoKey(orderId)
	pipe.Del(ctx, infoKey)

	_, err := pipe.Exec(ctx)
	return err
}
