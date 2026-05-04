package dao

import (
	"seckill/internal/paymentSvr/core/models"
)

func (r *Dao) GetOrderByOrderId(orderId string) (*models.Order, error) {
	var data models.Order

	err := r.Pgdb.Where("order_id = ?", orderId).First(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
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
