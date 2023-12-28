package models

import (
	"gorm.io/gorm"
)

type Accounts struct {
	gorm.Model
	Owner   string `gorm:"type:varchar(50);not null"`
	Balance uint64 `gorm:"type:bigInt;default:0;not null"`
}
