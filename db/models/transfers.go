package models

import (
	"gorm.io/gorm"
)

type Transfers struct {
	gorm.Model
	FromAccountID   uint64 `gorm:"type:bigint;not null"`
	ToAccountID     uint64 `gorm:"type:bigint;not null"`
	Amount          uint32 `gorm:"type:bigint;not null"`
	IncomingEntryID uint64 `gorm:"type:bigint;not null"`
	OutGoingEntryID uint64 `gorm:"type:bigint;not null"`
}
