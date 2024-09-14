package gorm

import "gorm.io/gorm"

type Product struct {
	ID           string `gorm:"primaryKey;column:id"`
	Name         string `gorm:"column:name"`
	Price        int    `gorm:"column:price"`
	CreatedAt    int64  `gorm:"column:created_at;autoCreateTime:mili"`
	UpdatedAt    int64  `gorm:"column:updated_at;autoCreateTime:mili;autoUpdateTime:mili"`
	DeletedAt    gorm.DeletedAt
	LikedByUsers []User `gorm:"many2many:user_like_product;foreignKey:id;joinForeignKey:product_id;references:id;joinReferences:user_id"`
}
