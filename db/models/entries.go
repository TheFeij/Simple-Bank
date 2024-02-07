package models

import (
	"gorm.io/gorm"
	"time"
)

type Entries struct {
	ID        int64          `gorm:"column:id"`
	AccountID int64          `gorm:"column:account_id"`
	Amount    int32          `gorm:"column:amount"`
	CreatedAt time.Time      `gorm:"column:created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at"`
}
