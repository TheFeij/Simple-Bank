package models

import (
	"gorm.io/gorm"
)

type Entries struct {
	gorm.Model
	AccountID uint64 `gorm:"column:account_id"`
	Amount    int64  `gorm:"column:amount"`
}
