package services

import (
	"Simple-Bank/requests"
	"Simple-Bank/responses"
	"Simple-Bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) responses.CreateAccountResponse {
	testAccount := requests.CreateAccountRequest{
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrency(),
		Balance:  uint64(util.RandomInt(0, 9999)),
	}

	account, err := accountServices.CreateAccount(testAccount)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, testAccount.Owner, account.Owner)
	require.Equal(t, testAccount.Currency, account.Currency)
	require.Equal(t, testAccount.Balance, account.Balance)

	require.NotZero(t, account.AccountID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)

	response, err := accountServices.GetAccount(uint64(account.AccountID))

	require.NoError(t, err)
	require.NotEmpty(t, response)

	require.Equal(t, account.AccountID, response.AccountID)
	require.Equal(t, account.Currency, response.Currency)
	require.Equal(t, account.Balance, response.Balance)
	//require.Equal(t, account.CreatedAt, response.CreatedAt)
	require.Equal(t, account.Owner, response.Owner)
	require.WithinDuration(t, account.CreatedAt, response.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	response, err := accountServices.DeleteAccount(account.AccountID)
	require.NoError(t, err)
	require.NotEmpty(t, response)

	response, err = accountServices.GetAccount(account.AccountID)
	require.Error(t, err)
	require.Empty(t, response)
}

func TestGetAccountsList(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	accounts, err := accountServices.ListAccounts(5)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
