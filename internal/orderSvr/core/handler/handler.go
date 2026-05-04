package handler

import (
	"context"
	"encoding/json"
	"seckill/internal/orderSvr/core/dto"
	"seckill/internal/orderSvr/kitex_gen/ordersvr"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

type OrderMessage struct {
	OrderId string  `json:"orderId"`
	UserId  string  `json:"userId"`
	ItemId  string  `json:"itemId"`
	Price   float64 `json:"price"`
}

func (s *OrderSvrImpl) CreateOrder(ctx context.Context, userId string, itemId string, price float64) (resp string, err error) {
	if len(userId) == 0 {
		return "", kerrors.NewBizStatusError(dto.InvalidUserID.Status, dto.InvalidUserID.Info)
	}
	if len(itemId) == 0 {
		return "", kerrors.NewBizStatusError(dto.InvalidItemID.Status, dto.InvalidItemID.Info)
	}
	if price <= 0 {
		return "", kerrors.NewBizStatusError(dto.InvalidPrice.Status, dto.InvalidPrice.Info)
	}

	orderId := uuid.New().String()

	msg := OrderMessage{
		OrderId: orderId,
		UserId:  userId,
		ItemId:  itemId,
		Price:   price,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}

	err = s.KafkaProd.WriteMessages(ctx, kafka.Message{
		Key:   []byte(orderId),
		Value: data,
	})
	if err != nil {
		return "", kerrors.NewBizStatusError(dto.MQSendFailed.Status, dto.MQSendFailed.Info)
	}

	return orderId, nil
}

func (s *OrderSvrImpl) queryOrdersByStatus(ctx context.Context, userId string, status int32) (resp []*ordersvr.OrderInfo, err error) {
	if len(userId) == 0 {
		return nil, kerrors.NewBizStatusError(dto.InvalidUserID.Status, dto.InvalidUserID.Info)
	}

	orders, err := s.Dao.GetOrdersByUserIdAndStatus(userId, status)
	if err != nil {
		return nil, err
	}

	var result []*ordersvr.OrderInfo
	for i := range orders {
		info := &ordersvr.OrderInfo{
			OrderId:    orders[i].OrderId,
			UserId:     orders[i].UserId,
			ItemId:     orders[i].ItemId,
			Price:      orders[i].Price,
			Status:     ordersvr.OrderStatus(orders[i].Status),
			CreateTime: orders[i].CreatedAt.Format("2006-01-02 15:04:05"),
		}
		result = append(result, info)
	}

	return result, nil
}

func (s *OrderSvrImpl) QueryPaidOrders(ctx context.Context, userId string) (resp []*ordersvr.OrderInfo, err error) {
	return s.queryOrdersByStatus(ctx, userId, 2)
}

func (s *OrderSvrImpl) QueryUnpaidOrders(ctx context.Context, userId string) (resp []*ordersvr.OrderInfo, err error) {
	return s.queryOrdersByStatus(ctx, userId, 1)
}

func (s *OrderSvrImpl) QueryCancelledOrders(ctx context.Context, userId string) (resp []*ordersvr.OrderInfo, err error) {
	return s.queryOrdersByStatus(ctx, userId, 3)
}
