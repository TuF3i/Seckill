package handler

import (
	"context"
	"seckill/internal/itemSvr/core/dto"
	"testing"

	"github.com/cloudwego/kitex/pkg/kerrors"
)

func newTestItemSvrImpl() *ItemSvrImpl {
	return NewItemSvrImpl(&ItemSvrImplReliance{
		Dao:   nil,
		Cache: nil,
	})
}

func TestAddItem_InvalidName(t *testing.T) {
	s := newTestItemSvrImpl()
	_, err := s.AddItem(context.Background(), "", 10, 9.99, "test")
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	bizErr, ok := kerrors.FromBizStatusError(err)
	if !ok {
		t.Fatalf("expected biz status error, got %v", err)
	}
	if bizErr.BizStatusCode() != dto.InvalidItemName.Status {
		t.Errorf("expected status %d, got %d", dto.InvalidItemName.Status, bizErr.BizStatusCode())
	}
}

func TestAddItem_InvalidStock(t *testing.T) {
	s := newTestItemSvrImpl()
	_, err := s.AddItem(context.Background(), "item1", 0, 9.99, "test")
	if err == nil {
		t.Fatal("expected error for zero stock")
	}
	bizErr, ok := kerrors.FromBizStatusError(err)
	if !ok {
		t.Fatalf("expected biz status error, got %v", err)
	}
	if bizErr.BizStatusCode() != dto.InvalidItemStock.Status {
		t.Errorf("expected status %d, got %d", dto.InvalidItemStock.Status, bizErr.BizStatusCode())
	}
}

func TestAddItem_InvalidPrice(t *testing.T) {
	s := newTestItemSvrImpl()
	_, err := s.AddItem(context.Background(), "item1", 10, 0, "test")
	if err == nil {
		t.Fatal("expected error for zero price")
	}
	bizErr, ok := kerrors.FromBizStatusError(err)
	if !ok {
		t.Fatalf("expected biz status error, got %v", err)
	}
	if bizErr.BizStatusCode() != dto.InvalidItemPrice.Status {
		t.Errorf("expected status %d, got %d", dto.InvalidItemPrice.Status, bizErr.BizStatusCode())
	}
}

func TestAddItem_LongName(t *testing.T) {
	s := newTestItemSvrImpl()
	longName := make([]byte, 129)
	for i := range longName {
		longName[i] = 'a'
	}
	_, err := s.AddItem(context.Background(), string(longName), 10, 9.99, "test")
	if err == nil {
		t.Fatal("expected error for too long name")
	}
	bizErr, ok := kerrors.FromBizStatusError(err)
	if !ok {
		t.Fatalf("expected biz status error, got %v", err)
	}
	if bizErr.BizStatusCode() != dto.InvalidItemName.Status {
		t.Errorf("expected status %d, got %d", dto.InvalidItemName.Status, bizErr.BizStatusCode())
	}
}

func TestDeleteItem_InvalidID(t *testing.T) {
	s := newTestItemSvrImpl()
	err := s.DeleteItem(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty id")
	}
	bizErr, ok := kerrors.FromBizStatusError(err)
	if !ok {
		t.Fatalf("expected biz status error, got %v", err)
	}
	if bizErr.BizStatusCode() != dto.InvalidItemID.Status {
		t.Errorf("expected status %d, got %d", dto.InvalidItemID.Status, bizErr.BizStatusCode())
	}
}

func TestStartFlashSale_InvalidID(t *testing.T) {
	s := newTestItemSvrImpl()
	err := s.StartFlashSale(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty id")
	}
	bizErr, ok := kerrors.FromBizStatusError(err)
	if !ok {
		t.Fatalf("expected biz status error, got %v", err)
	}
	if bizErr.BizStatusCode() != dto.InvalidItemID.Status {
		t.Errorf("expected status %d, got %d", dto.InvalidItemID.Status, bizErr.BizStatusCode())
	}
}

func TestStopFlashSale_InvalidID(t *testing.T) {
	s := newTestItemSvrImpl()
	err := s.StopFlashSale(context.Background(), "")
	if err == nil {
		t.Fatal("expected error for empty id")
	}
	bizErr, ok := kerrors.FromBizStatusError(err)
	if !ok {
		t.Fatalf("expected biz status error, got %v", err)
	}
	if bizErr.BizStatusCode() != dto.InvalidItemID.Status {
		t.Errorf("expected status %d, got %d", dto.InvalidItemID.Status, bizErr.BizStatusCode())
	}
}
