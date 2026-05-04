package handler

import (
	"context"
	"seckill/internal/paymentSvr/core/dto"
	"testing"

	"github.com/cloudwego/kitex/pkg/kerrors"
)

func newTestPaymentSvrImpl() *PaymentSvrImpl {
	return NewPaymentSvrImpl(&PaymentSvrImplReliance{
		Dao:   nil,
		Cache: nil,
	})
}

func TestProcessPayment_InvalidOrderID(t *testing.T) {
	s := newTestPaymentSvrImpl()
	_, err := s.ProcessPayment(context.Background(), "", "user1")
	if err == nil {
		t.Fatal("expected error for empty order id")
	}
	bizErr, ok := kerrors.FromBizStatusError(err)
	if !ok {
		t.Fatalf("expected biz status error, got %v", err)
	}
	if bizErr.BizStatusCode() != dto.InvalidOrderID.Status {
		t.Errorf("expected status %d, got %d", dto.InvalidOrderID.Status, bizErr.BizStatusCode())
	}
}
