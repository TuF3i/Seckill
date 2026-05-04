package dao

import (
	"seckill/internal/userSvr/core/models"
	"seckill/internal/userSvr/kitex_gen/usersvr"
	"seckill/pkg/enumTransfer"
)

func (r *Dao) AddUser(uid string, email string, password string) error {
	tx := r.pgdb.Begin()

	data := &models.User{
		Uid:      uid,
		Role:     enumTransfer.EnumToRoleString(usersvr.UserRole_SIMPLE_USER),
		Email:    email,
		Password: password,
	}

	err := tx.Create(data).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *Dao) GetUserInfo(email string) (*models.User, error) {
	var data models.User

	err := r.pgdb.Where("email = ?", email).First(&data).Error
	if err != nil {
		return nil, err
	}

	return &data, nil
}
