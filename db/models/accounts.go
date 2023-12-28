package models

import (
	"gorm.io/gorm"
)

type Accounts struct {
	gorm.Model
	Owner    string      `gorm:"type:varchar(50);not null"`
	Balance  uint64      `gorm:"type:bigInt;default:0;not null"`
	Entries  []Entries   `gorm:"foreignKey:AccountID;references:ID"`
	Incoming []Transfers `gorm:"foreignKey:ToAccountID;references:ID"`
	Outgoing []Transfers `gorm:"foreignKey:FromAccountID;references:ID"`
}
