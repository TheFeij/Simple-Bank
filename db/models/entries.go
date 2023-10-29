package models

import (
	"gorm.io/gorm"
)

type Entries struct {
	gorm.Model
	AccountID uint64 `gorm:"type:bigint;not null"`
	Amount    int64  `gorm:"type:bigint;not null"`
}
