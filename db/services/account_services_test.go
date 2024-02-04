package services

import (
	"Simple-Bank/requests"
	"Simple-Bank/responses"
	"Simple-Bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) responses.GetAccountResponse {
	testAccount := requests.CreateAccountRequest{
		Owner:   util.RandomOwner(),
		Balance: uint64(util.RandomInt(0, 9999)),
	}

	account, err := accountServices.CreateAccount(testAccount)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, testAccount.Owner, account.Owner)
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

	response, err := accountServices.GetAccount(account.AccountID)

	require.NoError(t, err)
	require.NotEmpty(t, response)

	require.Equal(t, account.AccountID, response.AccountID)
	require.Equal(t, account.Balance, response.Balance)
	require.Equal(t, account.Owner, response.Owner)
	require.WithinDuration(t, account.CreatedAt, response.CreatedAt, time.Second)
	require.Zero(t, response.DeletedAt)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	response, err := accountServices.DeleteAccount(account.AccountID)
	require.NoError(t, err)
	require.NotEmpty(t, response)

	response, err = accountServices.GetAccount(account.AccountID)
	require.NoError(t, err)
	require.NotEmpty(t, response)
	require.NotZero(t, response.DeletedAt)
}

func TestGetAccountsList(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	accounts, err := accountServices.ListAccounts(5)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.NotEmpty(t, accounts.Accounts)
	require.Len(t, accounts.Accounts, 5)

	for _, account := range accounts.Accounts {
		require.NotEmpty(t, account)
	}
}

func TestTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	concurrentTransactions := 5
	var amount int64 = 10

	errorsChan := make(chan error)
	resultsChan := make(chan responses.TransferResponse)

	for i := 0; i < concurrentTransactions; i++ {
		go func(chan responses.TransferResponse, chan error) {
			transferRequest := requests.TransferRequest{
				Amount:        uint32(amount),
				FromAccountID: account1.AccountID,
				ToAccountID:   account2.AccountID,
			}
			transfer, err := accountServices.Transfer(transferRequest)

			errorsChan <- err
			resultsChan <- transfer
		}(resultsChan, errorsChan)
	}

	for i := 0; i < concurrentTransactions; i++ {
		err := <-errorsChan
		require.NoError(t, err)

		transfer := <-resultsChan
		require.NotEmpty(t, transfer)
		require.Equal(t, amount, transfer.Amount)
		require.Equal(t, account1.AccountID, transfer.SrcAccountID)
		require.Equal(t, account2.AccountID, transfer.DstAccountID)
		require.NotZero(t, transfer.CreatedAt)
		require.NotZero(t, transfer.TransferID)

		var result responses.TransferResponse
		result, err = accountServices.GetTransfer(transfer.TransferID)
		require.NoError(t, err)

		var fromEntry responses.EntryResponse
		fromEntry, err = accountServices.GetEntry(result.OutgoingEntryID)
		require.NoError(t, err)
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.AccountID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.EntryID)
		require.NotZero(t, fromEntry.CreatedAt)

		var toEntry responses.EntryResponse
		toEntry, err = accountServices.GetEntry(result.IncomingEntryID)
		require.NoError(t, err)
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.AccountID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.EntryID)
		require.NotZero(t, toEntry.CreatedAt)
	}

	//TODO check account balances as well
}
