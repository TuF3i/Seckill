package handler

import (
	"context"
	"errors"
	"seckill/internal/itemSvr/core/dto"

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
		if errors.Is(err, redis.Nil) {
			return kerrors.NewBizStatusError(dto.ItemNotFound.Status, dto.ItemNotFound.Info)
		}
		return err
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
