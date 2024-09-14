package gorm

import "gorm.io/gorm"

type Address struct {
	ID        string `gorm:"primaryKey;column:id"`
	Address   string `gorm:"column:address"`
	UserId    string `gorm:"column:user_id"`
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:mili"`
	UpdatedAt int64  `gorm:"column:updated_at;autoCreateTime:mili;autoUpdateTime:mili"`
	DeletedAt gorm.DeletedAt
	User      User `gorm:"foreignKey:user_id;references:id"`
}

func (a *Address) TableName() string {
	return "addresses"
}
