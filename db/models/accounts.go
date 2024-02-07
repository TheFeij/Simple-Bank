package models

import (
	"gorm.io/gorm"
)

type Accounts struct {
	gorm.Model
	Owner   string `gorm:"column:owner"`
	Balance uint64 `gorm:"column:balance"`
}
