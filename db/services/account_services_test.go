package services

import (
	"Simple-Bank/db/models"
	"Simple-Bank/requests"
	"Simple-Bank/util"
	"github.com/stretchr/testify/require"
	"testing"
)

func createRandomAccount(t *testing.T) models.Accounts {
	testAccount := requests.CreateAccountRequest{
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrency(),
	}

	account, err := accountServices.CreateAccount(testAccount)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, testAccount.Owner, account.Owner)
	require.Equal(t, testAccount.Currency, account.Currency)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)

	response, err := accountServices.getAccount(uint64(account.ID))

	require.NoError(t, err)
	require.NotEmpty(t, response)

	require.Equal(t, account.ID, response.AccountID)
	require.Equal(t, account.Currency, response.Currency)
	require.Equal(t, account.Balance, response.Balance)
	require.Equal(t, account.CreatedAt, response.CreatedAt)
	require.Equal(t, account.Owner, response.Owner)
	//require.WithinDuration(t, account.CreatedAt, response.CreatedAt, time.Second)
}
