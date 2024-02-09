package models

import (
	"gorm.io/gorm"
	"time"
)

type Account struct {
	ID        int64          `gorm:"column:id"`
	Owner     string         `gorm:"column:owner"`
	Balance   int64          `gorm:"column:balance"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}
