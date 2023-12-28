package responses

import (
	"time"
)

type CreateAccountResponse struct {
	AccountID uint64    `json:"accountID"`
	CreatedAt time.Time `json:"createdAt"`
	Owner     string    `json:"owner"`
	Balance   uint64    `json:"balance"`
}

type GetAccountResponse struct {
	AccountID uint64    `json:"accountID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Owner     string    `json:"owner"`
	Balance   uint64    `json:"balance"`
}

type ListAccountsResponse struct {
	AccountID uint64    `json:"accountID"`
	CreatedAt time.Time `json:"createdAt"`
	UpdateAt  time.Time `json:"updatedAt"`
	Owner     string    `json:"owner"`
	Balance   uint64    `json:"balance"`
}

type DepositResponse struct {
	AccountID uint64    `json:"accountID"`
	Time      time.Time `json:"time"`
	Amount    int64     `json:"amount"`
}

type WithdrawResponse struct {
	AccountID uint64    `json:"accountID"`
	Time      time.Time `json:"time"`
	Amount    int64     `json:"amount"`
}
