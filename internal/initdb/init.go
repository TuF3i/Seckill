package initdb

import (
	"fmt"
	"log"
	n "seckill/infrastructures/nacos"
	"seckill/internal/userSvr/core/models"
	userSvr "seckill/internal/userSvr/kitex_gen/usersvr"
	"seckill/pkg/config"
	"seckill/pkg/enumTransfer"
	"seckill/pkg/env"
	"strconv"

	"seckill/infrastructures/postgres"
	itemModel "seckill/internal/itemSvr/core/models"
	orderModel "seckill/internal/orderSvr/core/models"

	"github.com/google/uuid"
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

	createAdminUser(db)

	log.Printf("[InitDB] database initialized successfully")
}

func autoMigrate(db *gorm.DB) error {
	models := []interface{}{
		&models.User{},
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

func createAdminUser(db *gorm.DB) {
	adminEmail := "admin@seckill.com"

	var count int64
	db.Model(&models.User{}).Where("email = ?", adminEmail).Count(&count)
	if count > 0 {
		log.Printf("[InitDB] admin user already exists (email=%s), skipping", adminEmail)
		return
	}

	uid := uuid.New().String()
	admin := &models.User{
		Uid:      uid,
		Role:     enumTransfer.EnumToRoleString(userSvr.UserRole_ADMIN),
		Email:    adminEmail,
		Password: "admin123",
	}

	if err := db.Create(admin).Error; err != nil {
		log.Printf("[InitDB] create admin user failed: %v", err)
		return
	}

	log.Printf("[InitDB] admin user created (email=%s, password=admin123, uid=%s)", adminEmail, uid)
}
