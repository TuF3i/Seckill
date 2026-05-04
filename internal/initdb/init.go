package initdb

import (
	"fmt"
	"log"

	"seckill/configs"
	"seckill/infrastructures/postgres"
	itemModel "seckill/internal/itemSvr/core/models"
	orderModel "seckill/internal/orderSvr/core/models"
	userModel "seckill/internal/userSvr/core/models"

	"gorm.io/gorm"
)

func InitDB(cfg *configs.PostgreSQLConfig) error {
	portStr := fmt.Sprintf("%d", cfg.Port)

	db, err := postgres.NewPostgresClient(
		postgres.WithHost(cfg.Host),
		postgres.WithPort(portStr),
		postgres.WithUser(cfg.User),
		postgres.WithPassword(cfg.Password),
		postgres.WithDefaultDB(cfg.DefaultDB),
		postgres.WithSSlMode(cfg.SSLMode),
	)
	if err != nil {
		return fmt.Errorf("initdb: connect to postgresql: %w", err)
	}

	err = autoMigrate(db)
	if err != nil {
		return fmt.Errorf("initdb: auto migrate: %w", err)
	}

	log.Printf("[InitDB] database initialized successfully")
	return nil
}

func autoMigrate(db *gorm.DB) error {
	models := []interface{}{
		&userModel.User{},
		&itemModel.Item{},
		&orderModel.Order{},
	}

	for _, m := range models {
		if err := db.AutoMigrate(m); err != nil {
			return fmt.Errorf("migrate %T: %w", m, err)
		}
	}

	return nil
}
