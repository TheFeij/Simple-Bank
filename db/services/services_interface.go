package services

import (
	"Simple-Bank/db/models"
	"Simple-Bank/requests"
)

type Services interface {
	CreateAccount(req requests.CreateAccountRequest) (models.Account, error)
	DeleteAccount(id int64) (models.Account, error)
	DepositMoney(req requests.DepositRequest) (models.Entry, error)
	WithdrawMoney(req requests.WithdrawRequest) (models.Entry, error)
	Transfer(req requests.TransferRequest) (models.Transfer, error)
	ListAccounts(pageNumber int64, pageSize int8) ([]models.Account, error)
	GetAccount(id int64) (models.Account, error)
	GetTransfer(id int64) (models.Transfer, error)
	GetEntry(id int64) (models.Entry, error)
	GetUser(username string) (models.User, error)
	CreateUser(req requests.CreateUserRequest) (models.User, error)
}

var _ Services = (*SQLServices)(nil)
