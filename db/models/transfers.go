package models

import (
	"gorm.io/gorm"
)

type Transfers struct {
	gorm.Model
	FromAccountID   uint64 `gorm:"column:from_account_id"`
	ToAccountID     uint64 `gorm:"column:to_account_id"`
	Amount          uint32 `gorm:"column:amount"`
	IncomingEntryID uint64 `gorm:"column:incoming_entry_id"`
	OutgoingEntryID uint64 `gorm:"column:outgoing_entry_id"`
}
