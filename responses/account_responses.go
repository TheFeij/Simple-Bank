package responses

import (
	"gorm.io/gorm"
	"time"
)

type CreateAccountResponse struct {
	AccountID int64     `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
	Owner     string    `json:"owner"`
	Balance   int64     `json:"balance"`
}

type GetAccountResponse struct {
	AccountID int64          `json:"account_id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
	Owner     string         `json:"owner"`
	Balance   int64          `json:"balance"`
}

type ListAccountsResponse struct {
	Accounts []GetAccountResponse `json:"accounts"`
}

type EntryResponse struct {
	EntryID   int64     `json:"entry_id"`
	AccountID int64     `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
	Amount    int32     `json:"amount"`
}
