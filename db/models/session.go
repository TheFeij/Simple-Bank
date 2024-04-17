package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Session struct {
	ID           uuid.UUID      `gorm:"column:id"`
	Username     string         `gorm:"column:username"`
	RefreshToken string         `gorm:"column:refresh_token"`
	UserAgent    string         `gorm:"column:user_agent"`
	ClientIP     string         `gorm:"column:client_ip"`
	IsBlocked    bool           `gorm:"column:is_blocked"`
	CreatedAt    time.Time      `gorm:"column:created_at"`
	ExpiresAt    time.Time      `gorm:"column:expires_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at"`
}
