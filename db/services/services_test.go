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
		Owner:   util.RandomUsername(),
		Balance: util.RandomBalance(),
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

	accounts, err := accountServices.ListAccounts(util.RandomInt(1, 2), 5)

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

	concurrentTransactions := 20
	var amount int32 = 10

	errorsChan := make(chan error)
	resultsChan := make(chan responses.TransferResponse)

	for i := 0; i < concurrentTransactions; i++ {
		go func(chan responses.TransferResponse, chan error) {
			transferRequest := requests.TransferRequest{
				Amount:        amount,
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
		require.Equal(t, amount, uint32(-fromEntry.Amount))
		require.NotZero(t, fromEntry.EntryID)
		require.NotZero(t, fromEntry.CreatedAt)

		var toEntry responses.EntryResponse
		toEntry, err = accountServices.GetEntry(result.IncomingEntryID)
		require.NoError(t, err)
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.AccountID, toEntry.AccountID)
		require.Equal(t, amount, uint32(toEntry.Amount))
		require.NotZero(t, toEntry.EntryID)
		require.NotZero(t, toEntry.CreatedAt)

		fromAccount, err := accountServices.GetAccount(account1.AccountID)
		require.Equal(t, account1.AccountID, fromAccount.AccountID)
		require.NoError(t, err)

		toAccount, err := accountServices.GetAccount(account2.AccountID)
		require.Equal(t, account2.AccountID, toAccount.AccountID)
		require.NoError(t, err)

		var diff1 = account1.Balance - fromAccount.Balance
		var diff2 = toAccount.Balance - account2.Balance

		require.True(t, diff2 > 0)
		require.True(t, diff1 > 0)
		require.Equal(t, diff1, diff2)
		require.Equal(t, uint64(amount), diff1/(int64(i)+int64(1)))
	}
}

func TestTransferDeadLock(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	concurrentTransactions := 20
	var amount int32 = 10

	errorsChan := make(chan error)
	resultsChan := make(chan responses.TransferResponse)

	for i := 0; i < concurrentTransactions; i++ {
		reverse := i%2 == 0
		go func(chan responses.TransferResponse, chan error, bool) {
			fromAccountID, toAccountID := account1.AccountID, account2.AccountID
			if reverse {
				fromAccountID, toAccountID = toAccountID, fromAccountID
			}
			transferRequest := requests.TransferRequest{
				Amount:        amount,
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
			}
			transfer, err := accountServices.Transfer(transferRequest)
			errorsChan <- err
			resultsChan <- transfer
		}(resultsChan, errorsChan, reverse)
	}

	for i := 0; i < concurrentTransactions; i++ {
		err := <-errorsChan
		require.NoError(t, err)

		transfer := <-resultsChan
		require.NotEmpty(t, transfer)

		var result responses.TransferResponse
		result, err = accountServices.GetTransfer(transfer.TransferID)
		require.NoError(t, err)
		require.NotEmpty(t, result)

		var fromEntry responses.EntryResponse
		fromEntry, err = accountServices.GetEntry(result.OutgoingEntryID)
		require.NoError(t, err)
		require.NotEmpty(t, fromEntry)

		var toEntry responses.EntryResponse
		toEntry, err = accountServices.GetEntry(result.IncomingEntryID)
		require.NoError(t, err)
		require.NotEmpty(t, toEntry)

		fromAccount, err := accountServices.GetAccount(result.SrcAccountID)
		require.NotEmpty(t, fromAccount)
		require.NoError(t, err)

		toAccount, err := accountServices.GetAccount(result.DstAccountID)
		require.NotEmpty(t, toAccount)
		require.NoError(t, err)
	}

	account1After, err := accountServices.GetAccount(account1.AccountID)
	require.NotEmpty(t, account1)
	require.NoError(t, err)
	require.Equal(t, account1.Balance, account1After.Balance)

	account2After, err := accountServices.GetAccount(account2.AccountID)
	require.NotEmpty(t, account2)
	require.NoError(t, err)
	require.Equal(t, account2.Balance, account2After.Balance)

}

func createRandomUser(t *testing.T) responses.UserInformationResponse {
	testUser := requests.CreateUserRequest{
		Username: util.RandomUsername(),
		Email:    util.RandomEmail(),
		FullName: util.RandomFullname(),
		Password: util.RandomPassword(),
	}

	createdTime := time.Now().Truncate(time.Nanosecond).Local()

	user, err := accountServices.CreateUser(testUser)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, testUser.Username, user.Username)
	require.Equal(t, testUser.Email, user.Email)
	//TODO check passwords
	require.Equal(t, testUser.FullName, user.FullName)
	require.WithinDuration(t, createdTime, user.CreatedAt, time.Second)
	require.WithinDuration(t, createdTime, user.UpdatedAt, time.Second)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}
