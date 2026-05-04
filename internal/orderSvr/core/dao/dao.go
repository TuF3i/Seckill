package dao

import (
	"seckill/internal/orderSvr/core/models"
)

func (r *Dao) GetOrderByOrderId(orderId string) (*models.Order, error) {
	var data models.Order

	err := r.Pgdb.Where("order_id = ?", orderId).First(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *Dao) GetOrdersByUserIdAndStatus(userId string, status int32) ([]models.Order, error) {
	var orders []models.Order

	err := r.Pgdb.Where("user_id = ? AND status = ?", userId, status).Find(&orders).Error
	if err != nil {
		return nil, err
	}

	return orders, nil
}
