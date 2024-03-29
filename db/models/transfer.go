package models

import (
	"gorm.io/gorm"
	"time"
)

type Transfer struct {
	ID              int64          `gorm:"column:id"`
	FromAccountID   int64          `gorm:"column:from_account_id"`
	ToAccountID     int64          `gorm:"column:to_account_id"`
	Amount          int32          `gorm:"column:amount"` // Amount range: [1, maxint32]
	IncomingEntryID int64          `gorm:"column:incoming_entry_id"`
	OutgoingEntryID int64          `gorm:"column:outgoing_entry_id"`
	CreatedAt       time.Time      `gorm:"column:created_at"`
	UpdatedAt       time.Time      `gorm:"column:updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"column:deleted_at"`
}
