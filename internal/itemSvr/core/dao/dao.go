package dao

import (
	"seckill/internal/itemSvr/core/models"
)

func (r *Dao) AddItem(itemId string, name string, stock int64, price float64, description string) error {
	tx := r.Pgdb.Begin()

	data := &models.Item{
		ItemId:      itemId,
		Name:        name,
		Stock:       stock,
		Price:       price,
		Description: description,
		FlashStatus: 0,
	}

	err := tx.Create(data).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Dao) DeleteItem(itemId string) error {
	tx := r.Pgdb.Begin()

	err := tx.Where("item_id = ?", itemId).Delete(&models.Item{}).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Dao) GetItem(itemId string) (*models.Item, error) {
	var data models.Item

	err := r.Pgdb.Where("item_id = ?", itemId).First(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *Dao) UpdateFlashStatus(itemId string, status int32) error {
	tx := r.Pgdb.Begin()

	err := tx.Model(&models.Item{}).Where("item_id = ?", itemId).Update("flash_status", status).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Dao) ListAllItems() ([]models.Item, error) {
	var items []models.Item

	err := r.Pgdb.Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Dao) HasActiveFlashSale() (bool, error) {
	var count int64

	err := r.Pgdb.Model(&models.Item{}).Where("flash_status = 1").Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
