package handler

import (
	"context"
	"seckill/internal/paymentSvr/core/dto"

	"github.com/cloudwego/kitex/pkg/kerrors"
)

func (s *PaymentSvrImpl) ProcessPayment(ctx context.Context, orderId string, userId string) (resp bool, err error) {
	if len(orderId) == 0 {
		return false, kerrors.NewBizStatusError(dto.InvalidOrderID.Status, dto.InvalidOrderID.Info)
	}

	order, err := s.Dao.GetOrderByOrderId(orderId)
	if err != nil {
		return false, kerrors.NewBizStatusError(dto.OrderNotFound.Status, dto.OrderNotFound.Info)
	}

	if order.Status == 2 {
		return false, kerrors.NewBizStatusError(dto.OrderAlreadyPaid.Status, dto.OrderAlreadyPaid.Info)
	}

	err = s.Dao.UpdateOrderStatus(orderId, 2)
	if err != nil {
		return false, kerrors.NewBizStatusError(dto.PaymentFailed.Status, dto.PaymentFailed.Info)
	}

	err = s.Cache.DelOrderCache(ctx, orderId)
	if err != nil {
		return false, kerrors.NewBizStatusError(dto.PaymentFailed.Status, dto.PaymentFailed.Info)
	}

	return true, nil
}
