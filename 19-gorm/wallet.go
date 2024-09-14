package gorm

import "gorm.io/gorm"

type Wallet struct {
	ID        string `gorm:"primaryKey;column:id"`
	Balance   int64  `gorm:"column:balance"`
	UserId    string `gorm:"column:user_id"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:mili"`
	UpdatedAt int64  `gorm:"column:updated_at;autoCreateTime:mili;autoUpdateTime:mili"`
	DeletedAt gorm.DeletedAt
	User      *User `gorm:"foreignKey:user_id;references:id"`
}
