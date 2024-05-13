package services

import (
	"Simple-Bank/db/models"
	"Simple-Bank/requests"
	"github.com/google/uuid"
)

type Services interface {
	CreateAccount(owner string) (models.Account, error)
	DeleteAccount(id int64) (models.Account, error)
	DepositMoney(req requests.DepositRequest) (models.Entry, error)
	WithdrawMoney(req requests.WithdrawRequest) (models.Entry, error)
	Transfer(req TransferRequest) (models.Transfer, error)
	ListAccounts(req ListAccountsRequest) ([]models.Account, error)
	GetAccount(id int64) (models.Account, error)
	GetTransfer(id int64) (models.Transfer, error)
	GetEntry(id int64) (models.Entry, error)
	GetUser(username string) (models.User, error)
	CreateUser(req requests.CreateUserRequest) (models.User, error)
	GetSession(id uuid.UUID) (models.Session, error)
	CreateSession(session models.Session) (models.Session, error)
	UpdateUser(req UpdateUserRequest) (models.User, error)
}

var _ Services = (*SQLServices)(nil)
