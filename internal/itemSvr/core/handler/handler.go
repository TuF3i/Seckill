package handler

import (
	"context"
	"errors"
	"seckill/internal/itemSvr/core/dto"
	itemsvr "seckill/internal/itemSvr/kitex_gen/itemsvr"
	"time"

	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

func (s *ItemSvrImpl) AddItem(ctx context.Context, name string, stock int64, price float64, description string) (resp string, err error) {
	if len(name) == 0 || len(name) > 128 {
		return "", kerrors.NewBizStatusError(dto.InvalidItemName.Status, dto.InvalidItemName.Info)
	}
	if stock <= 0 {
		return "", kerrors.NewBizStatusError(dto.InvalidItemStock.Status, dto.InvalidItemStock.Info)
	}
	if price <= 0 {
		return "", kerrors.NewBizStatusError(dto.InvalidItemPrice.Status, dto.InvalidItemPrice.Info)
	}

	active, err := s.Dao.HasActiveFlashSale()
	if err == nil && active {
		return "", kerrors.NewBizStatusError(dto.FlashSaleActive.Status, dto.FlashSaleActive.Info)
	}

	itemId := uuid.New().String()

	err = s.Dao.AddItem(itemId, name, stock, price, description)
	if err != nil {
		return "", err
	}

	return itemId, nil
}

func (s *ItemSvrImpl) DeleteItem(ctx context.Context, id string) (err error) {
	if len(id) == 0 {
		return kerrors.NewBizStatusError(dto.InvalidItemID.Status, dto.InvalidItemID.Info)
	}

	_, err = s.Dao.GetItem(id)
	if err != nil {
		return kerrors.NewBizStatusError(dto.ItemNotFound.Status, dto.ItemNotFound.Info)
	}

	active, err := s.Dao.HasActiveFlashSale()
	if err == nil && active {
		return kerrors.NewBizStatusError(dto.FlashSaleActive.Status, dto.FlashSaleActive.Info)
	}

	err = s.Dao.DeleteItem(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *ItemSvrImpl) StartFlashSale(ctx context.Context, itemId string) (err error) {
	if len(itemId) == 0 {
		return kerrors.NewBizStatusError(dto.InvalidItemID.Status, dto.InvalidItemID.Info)
	}

	item, err := s.Dao.GetItem(itemId)
	if err != nil {
		return kerrors.NewBizStatusError(dto.ItemNotFound.Status, dto.ItemNotFound.Info)
	}

	flashStatus, err := s.Cache.GetFlashStatus(ctx, itemId)
	if err == nil && flashStatus == 1 {
		return kerrors.NewBizStatusError(dto.FlashAlreadyStart.Status, dto.FlashAlreadyStart.Info)
	}

	err = s.Dao.UpdateFlashStatus(itemId, 1)
	if err != nil {
		return err
	}

	err = s.Cache.WarmUpItemStock(ctx, itemId, item.Stock)
	if err != nil {
		return err
	}

	err = s.Cache.SetFlashStatus(ctx, itemId, 1)
	if err != nil {
		return err
	}

	return nil
}

func (s *ItemSvrImpl) StopFlashSale(ctx context.Context, itemId string) (err error) {
	if len(itemId) == 0 {
		return kerrors.NewBizStatusError(dto.InvalidItemID.Status, dto.InvalidItemID.Info)
	}

	flashStatus, err := s.Cache.GetFlashStatus(ctx, itemId)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return kerrors.NewBizStatusError(dto.FlashNotStart.Status, dto.FlashNotStart.Info)
		}
		return err
	}
	if flashStatus == 0 {
		return kerrors.NewBizStatusError(dto.FlashNotStart.Status, dto.FlashNotStart.Info)
	}

	err = s.Dao.UpdateFlashStatus(itemId, 0)
	if err != nil {
		return err
	}

	err = s.Cache.DelItemFlashCache(ctx, itemId)
	if err != nil {
		return err
	}

	return nil
}

func (s *ItemSvrImpl) ListItems(ctx context.Context, uid string, role string) (resp []*itemsvr.ItemInfo, err error) {
	active, err := s.Dao.HasActiveFlashSale()
	if err != nil {
		return nil, err
	}

	if !active {
		if role != "ADMIN" {
			return nil, kerrors.NewBizStatusError(dto.FlashNotStart.Status, dto.FlashNotStart.Info)
		}

		items, err := s.Dao.ListAllItems()
		if err != nil {
			return nil, err
		}

		var result []*itemsvr.ItemInfo
		for i := range items {
			result = append(result, &itemsvr.ItemInfo{
				ID:          items[i].ItemId,
				Name:        items[i].Name,
				Stock:       items[i].Stock,
				Price:       items[i].Price,
				Description: items[i].Description,
			})
		}
		return result, nil
	}

	items, err := s.Dao.ListAllItems()
	if err != nil {
		return nil, err
	}

	var result []*itemsvr.ItemInfo
	for i := range items {
		info := &itemsvr.ItemInfo{
			ID:          items[i].ItemId,
			Name:        items[i].Name,
			Stock:       items[i].Stock,
			Price:       items[i].Price,
			Description: items[i].Description,
		}

		flashStatus, cacheErr := s.Cache.GetFlashStatus(ctx, items[i].ItemId)
		if cacheErr == nil && flashStatus == 1 {
			stock, cacheErr := s.Cache.GetItemStock(ctx, items[i].ItemId)
			if cacheErr == nil {
				info.Stock = stock
			}
		}

		result = append(result, info)
	}

	return result, nil
}

func (s *ItemSvrImpl) GetItem(ctx context.Context, itemId string) (resp *itemsvr.ItemInfo, err error) {
	if len(itemId) == 0 {
		return nil, kerrors.NewBizStatusError(dto.InvalidItemID.Status, dto.InvalidItemID.Info)
	}

	item, err := s.Dao.GetItem(itemId)
	if err != nil {
		return nil, kerrors.NewBizStatusError(dto.ItemNotFound.Status, dto.ItemNotFound.Info)
	}

	result := &itemsvr.ItemInfo{
		ID:          item.ItemId,
		Name:        item.Name,
		Stock:       item.Stock,
		Price:       item.Price,
		Description: item.Description,
	}

	flashStatus, cacheErr := s.Cache.GetFlashStatus(ctx, itemId)
	if cacheErr == nil && flashStatus == 1 {
		stock, cacheErr := s.Cache.GetItemStock(ctx, itemId)
		if cacheErr == nil {
			result.Stock = stock
		}
	}

	return result, nil
}

func (s *ItemSvrImpl) PrepareOrder(ctx context.Context, userId string, itemId string) (resp float64, err error) {
	if len(userId) == 0 {
		return 0, kerrors.NewBizStatusError(dto.InvalidItemID.Status, dto.InvalidItemID.Info)
	}
	if len(itemId) == 0 {
		return 0, kerrors.NewBizStatusError(dto.InvalidItemID.Status, dto.InvalidItemID.Info)
	}

	var result int64
	if s.Benchmark {
		result, err = s.Cache.PrepareOrderAtomicNoLimit(ctx, itemId)
	} else {
		result, err = s.Cache.PrepareOrderAtomic(ctx, itemId, userId, 5*time.Minute)
	}
	if err != nil {
		return 0, err
	}

	switch {
	case result == -1:
		return 0, kerrors.NewBizStatusError(dto.FlashNotStart.Status, dto.FlashNotStart.Info)
	case result == -2:
		return 0, kerrors.NewBizStatusError(dto.PurchaseLimitExceeded.Status, dto.PurchaseLimitExceeded.Info)
	case result == -3:
		return 0, kerrors.NewBizStatusError(dto.InsufficientStock.Status, dto.InsufficientStock.Info)
	}

	item, err := s.Dao.GetItem(itemId)
	if err != nil {
		return 0, kerrors.NewBizStatusError(dto.ItemNotFound.Status, dto.ItemNotFound.Info)
	}

	return item.Price, nil
}
