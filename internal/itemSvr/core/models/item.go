package models

import (
	"time"

	"gorm.io/gorm"
)

type Item struct {
	ID        int64          `gorm:"primaryKey;type:bigint;autoIncrement;comment:表自增主键ID"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;index;comment:删除时间（软删除）"`

	ItemId      string  `json:"itemId" gorm:"column:item_id;type:varchar(64);not null;uniqueIndex;comment:商品唯一ID"`
	Name        string  `json:"name" gorm:"column:name;type:varchar(128);not null;comment:商品名称"`
	Stock       int64   `json:"stock" gorm:"column:stock;type:bigint;not null;default:0;comment:商品库存"`
	Price       float64 `json:"price" gorm:"column:price;type:decimal(10,2);not null;default:0;comment:商品价格"`
	Description string  `json:"description" gorm:"column:description;type:text;comment:商品描述"`
	FlashStatus int32   `json:"flashStatus" gorm:"column:flash_status;type:int;not null;default:0;comment:秒杀状态 0-未开始 1-进行中"`
}

func (Item) TableName() string {
	return "item_table"
}
