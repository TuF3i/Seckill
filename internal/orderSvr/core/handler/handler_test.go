package handler

import (
	"context"
	"seckill/internal/orderSvr/core/dto"
	"testing"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/segmentio/kafka-go"
)

func newTestOrderSvrImpl() *OrderSvrImpl {
	return NewOrderSvrImpl(&OrderSvrImplReliance{
		Dao:       nil,
		Cache:     nil,
		KafkaProd: &kafka.Writer{},
	})
}

func TestCreateOrder_InvalidUserID(t *testing.T) {
	s := newTestOrderSvrImpl()
	_, err := s.CreateOrder(context.Background(), "", "item1", 9.99)
	if err == nil {
		t.Fatal("expected error for empty user id")
	}
	bizErr, ok := kerrors.FromBizStatusError(err)
	if !ok {
		t.Fatalf("expected biz status error, got %v", err)
	}
	if bizErr.BizStatusCode() != dto.InvalidUserID.Status {
		t.Errorf("expected status %d, got %d", dto.InvalidUserID.Status, bizErr.BizStatusCode())
	}
}

func TestCreateOrder_InvalidItemID(t *testing.T) {
	s := newTestOrderSvrImpl()
	_, err := s.CreateOrder(context.Background(), "user1", "", 9.99)
	if err == nil {
		t.Fatal("expected error for empty item id")
	}
	bizErr, ok := kerrors.FromBizStatusError(err)
	if !ok {
		t.Fatalf("expected biz status error, got %v", err)
	}
	if bizErr.BizStatusCode() != dto.InvalidItemID.Status {
		t.Errorf("expected status %d, got %d", dto.InvalidItemID.Status, bizErr.BizStatusCode())
	}
}

func TestCreateOrder_InvalidPrice(t *testing.T) {
	s := newTestOrderSvrImpl()
	_, err := s.CreateOrder(context.Background(), "user1", "item1", 0)
	if err == nil {
		t.Fatal("expected error for zero price")
	}
	bizErr, ok := kerrors.FromBizStatusError(err)
	if !ok {
		t.Fatalf("expected biz status error, got %v", err)
	}
	if bizErr.BizStatusCode() != dto.InvalidPrice.Status {
		t.Errorf("expected status %d, got %d", dto.InvalidPrice.Status, bizErr.BizStatusCode())
	}
}

func TestQueryPaidOrders_InvalidUserID(t *testing.T) {
	s := newTestOrderSvrImpl()
	_, err := s.QueryPaidOrders(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty user id")
	}
	bizErr, ok := kerrors.FromBizStatusError(err)
	if !ok {
		t.Fatalf("expected biz status error, got %v", err)
	}
	if bizErr.BizStatusCode() != dto.InvalidUserID.Status {
		t.Errorf("expected status %d, got %d", dto.InvalidUserID.Status, bizErr.BizStatusCode())
	}
}

func TestQueryUnpaidOrders_InvalidUserID(t *testing.T) {
	s := newTestOrderSvrImpl()
	_, err := s.QueryUnpaidOrders(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty user id")
	}
	bizErr, ok := kerrors.FromBizStatusError(err)
	if !ok {
		t.Fatalf("expected biz status error, got %v", err)
	}
	if bizErr.BizStatusCode() != dto.InvalidUserID.Status {
		t.Errorf("expected status %d, got %d", dto.InvalidUserID.Status, bizErr.BizStatusCode())
	}
}

func TestQueryCancelledOrders_InvalidUserID(t *testing.T) {
	s := newTestOrderSvrImpl()
	_, err := s.QueryCancelledOrders(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty user id")
	}
	bizErr, ok := kerrors.FromBizStatusError(err)
	if !ok {
		t.Fatalf("expected biz status error, got %v", err)
	}
	if bizErr.BizStatusCode() != dto.InvalidUserID.Status {
		t.Errorf("expected status %d, got %d", dto.InvalidUserID.Status, bizErr.BizStatusCode())
	}
}
