package services

import "fmt"

var (
	ErrUnAuthorizedWithdraw = fmt.Errorf("cannot withdraw money from other user's accounts")
	ErrUnAuthorizedDeposit  = fmt.Errorf("cannot deposit money into other user's accounts")
	ErrUserNotFound         = fmt.Errorf("user not found")
	ErrAccountNotFound      = fmt.Errorf("account not found")
)
