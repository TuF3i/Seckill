package initdb

import (
	"fmt"
	"log"
	n "seckill/infrastructures/nacos"
	"seckill/pkg/config"
	"seckill/pkg/env"
	"strconv"

	"seckill/infrastructures/postgres"
	itemModel "seckill/internal/itemSvr/core/models"
	orderModel "seckill/internal/orderSvr/core/models"
	userModel "seckill/internal/userSvr/core/models"

	"gorm.io/gorm"
)

func InitDB() {
	fmt.Println("Initializing database...")

	basicEnv := env.GetEnv()

	nacosPort, err := strconv.ParseUint(basicEnv.NacosPort, 10, 64)
	if err != nil {
		panic(fmt.Errorf("initdb: parse nacos port: %w", err))
	}

	nacosClient, err := n.NewNacosClient(
		n.WithHost(basicEnv.NacosAddr),
		n.WithPort(nacosPort),
		n.WithUserName(basicEnv.NacosUser),
		n.WithPassword(basicEnv.NacosPassword),
	)
	if err != nil {
		panic(fmt.Errorf("initdb: create nacos client: %w", err))
	}

	loader, err := config.NewLoader(nacosClient, basicEnv.ConfigID, basicEnv.ConfigGroup)
	if err != nil {
		panic(fmt.Errorf("initdb: load config: %w", err))
	}

	cfg := loader.GetConfig().PostgreSQL

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
		panic(fmt.Errorf("initdb: connect to postgresql: %w", err))
	}

	err = autoMigrate(db)
	if err != nil {
		panic(fmt.Errorf("initdb: auto migrate: %w", err))
	}

	log.Printf("[InitDB] database initialized successfully")
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
