package responses

import (
	"time"
)

type CreateAccountResponse struct {
	AccountID uint64    `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
	Owner     string    `json:"owner"`
	Balance   uint64    `json:"balance"`
}

type GetAccountResponse struct {
	AccountID uint64    `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Owner     string    `json:"owner"`
	Balance   uint64    `json:"balance"`
}

type ListAccountsResponse struct {
	AccountID uint64    `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdateAt  time.Time `json:"update_at"`
	Owner     string    `json:"owner"`
	Balance   uint64    `json:"balance"`
}

type DepositResponse struct {
	AccountID uint64    `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
	Amount    int64     `json:"amount"`
}

type WithdrawResponse struct {
	AccountID uint64    `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
	Amount    int64     `json:"amount"`
}
