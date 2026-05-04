package cache

import (
	"context"
	"encoding/json"
	"seckill/internal/orderSvr/core/models"
	"seckill/internal/orderSvr/core/pkg/lkeygen"
)

func (r *Cache) SetOrderStatus(ctx context.Context, orderId string, status int32) error {
	key := lkeygen.GenOrderStatusKey(orderId)
	return r.Rdb.Set(ctx, key, status, 0).Err()
}

func (r *Cache) GetOrderStatus(ctx context.Context, orderId string) (int32, error) {
	key := lkeygen.GenOrderStatusKey(orderId)
	val, err := r.Rdb.Get(ctx, key).Int64()
	if err != nil {
		return 0, err
	}
	return int32(val), nil
}

func (r *Cache) CacheOrderInfo(ctx context.Context, order *models.Order) error {
	key := lkeygen.GenOrderInfoKey(order.OrderId)
	data, err := json.Marshal(order)
	if err != nil {
		return err
	}
	return r.Rdb.Set(ctx, key, string(data), 0).Err()
}

func (r *Cache) GetOrderInfo(ctx context.Context, orderId string) (*models.Order, error) {
	key := lkeygen.GenOrderInfoKey(orderId)
	data, err := r.Rdb.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	var order models.Order
	err = json.Unmarshal([]byte(data), &order)
	if err != nil {
		return nil, err
	}
	return &order, nil
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

func (r *Cache) AddUserOrder(ctx context.Context, userId string, orderId string) error {
	key := lkeygen.GenUserOrdersKey(userId)
	return r.Rdb.SAdd(ctx, key, orderId).Err()
}
