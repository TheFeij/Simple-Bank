package services

import "fmt"

var (
	ErrUnAuthorizedWithdraw = fmt.Errorf("cannot withdraw money from other user's accounts")
	ErrUnAuthorizedDeposit  = fmt.Errorf("cannot deposit money into other user's accounts")
	ErrUnAuthorizedTransfer = fmt.Errorf("cannot transfer money from other user's accounts")
	ErrUserNotFound         = fmt.Errorf("user not found")
	ErrSrcAccountNotFound   = fmt.Errorf("source account not found")
	ErrDstAccountNotFound   = fmt.Errorf("destination account not found")
	ErrAccountNotFound      = fmt.Errorf("account not found")
	ErrTransferNotFound     = fmt.Errorf("transfer not found")
	ErrSessionNotFound      = fmt.Errorf("session not found")
	ErrEntryNotFound        = fmt.Errorf("entry not found")
	ErrNotEnoughMoney       = fmt.Errorf("not enough money in the account")
)
