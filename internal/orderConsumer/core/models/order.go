package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	ID        int64          `gorm:"primaryKey;type:bigint;autoIncrement;comment:表自增主键ID"`
	CreatedAt time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp;index;comment:删除时间（软删除）"`

	OrderId string  `json:"orderId" gorm:"column:order_id;type:varchar(64);not null;uniqueIndex;comment:订单唯一ID"`
	UserId  string  `json:"userId" gorm:"column:user_id;type:varchar(64);not null;index;comment:用户ID"`
	ItemId  string  `json:"itemId" gorm:"column:item_id;type:varchar(64);not null;comment:商品ID"`
	Price   float64 `json:"price" gorm:"column:price;type:decimal(10,2);not null;default:0;comment:订单金额"`
	Status  int32   `json:"status" gorm:"column:status;type:int;not null;default:1;comment:订单状态 1-未支付 2-已支付 3-已取消"`
}

func (Order) TableName() string {
	return "order_table"
}
