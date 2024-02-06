package services

import (
	"Simple-Bank/requests"
	"Simple-Bank/responses"
)

type Services interface {
	CreateAccount(req requests.CreateAccountRequest) (responses.CreateAccountResponse, error)
	DeleteAccount(id uint64) (responses.GetAccountResponse, error)
	DepositMoney(req requests.DepositRequest) (responses.EntryResponse, error)
	WithdrawMoney(req requests.WithdrawRequest) (responses.EntryResponse, error)
	Transfer(req requests.TransferRequest) (responses.TransferResponse, error)
	ListAccounts(pageNumber, pageSize uint64) (responses.ListAccountsResponse, error)
	GetAccount(id uint64) (responses.GetAccountResponse, error)
	GetTransfer(id uint64) (responses.TransferResponse, error)
	GetEntry(id uint64) (responses.EntryResponse, error)
}

var _ Services = (*SQLServices)(nil)
