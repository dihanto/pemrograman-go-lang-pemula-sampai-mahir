package gorm

import (
	"gorm.io/gorm"
)

type Todo struct {
	ID          int64          `gorm:"primaryKey;column:id;autoIncrement"`
	UserId      string         `gorm:"column:user_id"`
	Title       string         `gorm:"column:title"`
	Description string         `gorm:"column:description"`
	CreatedAt   []uint8        `gorm:"column:created_at;autoCreateTime:mili"`
	UpdatedAt   []uint8        `gorm:"column:updated_at;autoCreateTime:mili;autoUpdateTime:mili"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at"`
}
