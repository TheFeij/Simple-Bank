package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Username       string         `gorm:"column:username"`
	HashedPassword string         `gorm:"column:hashed_password"`
	FullName       string         `gorm:"column:fullname"`
	Email          string         `gorm:"column:email"`
	CreatedAt      time.Time      `gorm:"column:created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at"`
}
