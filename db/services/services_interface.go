package services

import (
	"Simple-Bank/db/models"
	"Simple-Bank/requests"
	"Simple-Bank/responses"
)

type Services interface {
	CreateAccount(req requests.CreateAccountRequest) (responses.CreateAccountResponse, error)
	DeleteAccount(id int64) (responses.GetAccountResponse, error)
	DepositMoney(req requests.DepositRequest) (responses.EntryResponse, error)
	WithdrawMoney(req requests.WithdrawRequest) (responses.EntryResponse, error)
	Transfer(req requests.TransferRequest) (responses.TransferResponse, error)
	ListAccounts(pageNumber int64, pageSize int8) (responses.ListAccountsResponse, error)
	GetAccount(id int64) (responses.GetAccountResponse, error)
	GetTransfer(id int64) (responses.TransferResponse, error)
	GetEntry(id int64) (responses.EntryResponse, error)
	GetUser(username string) (models.User, error)
	CreateUser(req requests.CreateUserRequest) (models.User, error)
}

var _ Services = (*SQLServices)(nil)
