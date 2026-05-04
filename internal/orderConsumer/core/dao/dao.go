package dao

import (
	"seckill/internal/orderConsumer/core/models"
)

func (r *Dao) SaveOrder(order *models.Order) error {
	tx := r.Pgdb.Begin()

	err := tx.Create(order).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Dao) UpdateOrderStatus(orderId string, status int32) error {
	tx := r.Pgdb.Begin()

	err := tx.Model(&models.Order{}).Where("order_id = ?", orderId).Update("status", status).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Dao) GetUnpaidOrders() ([]models.Order, error) {
	var orders []models.Order

	err := r.Pgdb.Where("status = 1").Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}
