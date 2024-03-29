package services

import (
	"Simple-Bank/db/models"
	"Simple-Bank/requests"
	"Simple-Bank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func createRandomAccount(t *testing.T) models.Account {
	user := createRandomUser(t)

	createAccountRequest := requests.CreateAccountRequest{
		Owner:   user.Username,
		Balance: util.RandomBalance(),
	}

	createdTime := time.Now().Truncate(time.Nanosecond).Local()

	account, err := services.CreateAccount(createAccountRequest)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, createAccountRequest.Owner, account.Owner)
	require.Equal(t, createAccountRequest.Balance, account.Balance)
	require.True(t, account.ID > 0)
	require.WithinDuration(t, createdTime, account.CreatedAt, time.Second)
	require.WithinDuration(t, createdTime, account.UpdatedAt, time.Second)
	require.True(t, account.DeletedAt.Time.IsZero())

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)

	response, err := services.GetAccount(account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, response)

	require.Equal(t, account.ID, response.ID)
	require.Equal(t, account.Balance, response.Balance)
	require.Equal(t, account.Owner, response.Owner)
	require.WithinDuration(t, account.CreatedAt, response.CreatedAt, time.Second)
	require.WithinDuration(t, account.UpdatedAt, response.UpdatedAt, time.Second)
	require.Equal(t, account.DeletedAt, response.DeletedAt)
	require.True(t, response.DeletedAt.Time.IsZero())
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	response, err := services.DeleteAccount(account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, response)
	require.Equal(t, account.ID, response.ID)
	require.Equal(t, account.Balance, response.Balance)
	require.Equal(t, account.Owner, response.Owner)
	require.WithinDuration(t, account.CreatedAt, response.CreatedAt, time.Second)
	require.WithinDuration(t, account.UpdatedAt, response.UpdatedAt, time.Second)
	require.NotEqual(t, account.DeletedAt, response.DeletedAt)
	require.False(t, response.DeletedAt.Time.IsZero())

	response, err = services.GetAccount(account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, response)
	require.False(t, response.DeletedAt.Time.IsZero())
}

func TestGetAccountsList(t *testing.T) {
	createdAccounts := make([]models.Account, 5)
	for i := 0; i < 5; i++ {
		createdAccounts[i] = createRandomAccount(t)
	}

	accounts, err := services.ListAccounts(1, 5)

	require.NoError(t, err)
	require.NotEmpty(t, accounts)
	require.Len(t, accounts, 5)

	for i, account := range accounts {
		require.NotEmpty(t, account)

		require.Equal(t, createdAccounts[i].Owner, account.Owner)
		require.Equal(t, createdAccounts[i].Balance, account.Balance)
		require.True(t, account.ID > 0)
		require.WithinDuration(t, createdAccounts[i].CreatedAt, account.CreatedAt, time.Second)
		require.WithinDuration(t, createdAccounts[i].UpdatedAt, account.UpdatedAt, time.Second)
		require.True(t, account.DeletedAt.Time.IsZero())
	}
}

func TestTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	concurrentTransactions := 20
	var amount int32 = 10

	errorsChan := make(chan error)
	resultsChan := make(chan models.Transfer)

	for i := 0; i < concurrentTransactions; i++ {
		go func(chan models.Transfer, chan error) {
			transferRequest := requests.TransferRequest{
				Amount:        amount,
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
			}
			transfer, err := services.Transfer(transferRequest)

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
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.NotZero(t, transfer.CreatedAt)
		require.NotZero(t, transfer.ID)

		var result models.Transfer
		result, err = services.GetTransfer(transfer.ID)
		require.NoError(t, err)

		var fromEntry models.Entry
		fromEntry, err = services.GetEntry(result.OutgoingEntryID)
		require.NoError(t, err)
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, amount, -fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		var toEntry models.Entry
		toEntry, err = services.GetEntry(result.IncomingEntryID)
		require.NoError(t, err)
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		fromAccount, err := services.GetAccount(account1.ID)
		require.Equal(t, account1.ID, fromAccount.ID)
		require.NoError(t, err)

		toAccount, err := services.GetAccount(account2.ID)
		require.Equal(t, account2.ID, toAccount.ID)
		require.NoError(t, err)

		var diff1 = int32(account1.Balance - fromAccount.Balance)
		var diff2 = int32(toAccount.Balance - account2.Balance)

		require.True(t, diff2 > 0)
		require.True(t, diff1 > 0)
		require.Equal(t, diff1, diff2)
		require.Equal(t, amount, diff1/int32(i+1))
	}
}

func TestTransferDeadLock(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	concurrentTransactions := 20
	var amount int32 = 10

	errorsChan := make(chan error)
	resultsChan := make(chan models.Transfer)

	for i := 0; i < concurrentTransactions; i++ {
		reverse := i%2 == 0
		go func(chan models.Transfer, chan error, bool) {
			fromAccountID, toAccountID := account1.ID, account2.ID
			if reverse {
				fromAccountID, toAccountID = toAccountID, fromAccountID
			}
			transferRequest := requests.TransferRequest{
				Amount:        amount,
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
			}
			transfer, err := services.Transfer(transferRequest)
			errorsChan <- err
			resultsChan <- transfer
		}(resultsChan, errorsChan, reverse)
	}

	for i := 0; i < concurrentTransactions; i++ {
		err := <-errorsChan
		require.NoError(t, err)

		transfer := <-resultsChan
		require.NotEmpty(t, transfer)

		var result models.Transfer
		result, err = services.GetTransfer(transfer.ID)
		require.NoError(t, err)
		require.NotEmpty(t, result)

		var fromEntry models.Entry
		fromEntry, err = services.GetEntry(result.OutgoingEntryID)
		require.NoError(t, err)
		require.NotEmpty(t, fromEntry)

		var toEntry models.Entry
		toEntry, err = services.GetEntry(result.IncomingEntryID)
		require.NoError(t, err)
		require.NotEmpty(t, toEntry)

		fromAccount, err := services.GetAccount(result.FromAccountID)
		require.NotEmpty(t, fromAccount)
		require.NoError(t, err)

		toAccount, err := services.GetAccount(result.ToAccountID)
		require.NotEmpty(t, toAccount)
		require.NoError(t, err)
	}

	account1After, err := services.GetAccount(account1.ID)
	require.NotEmpty(t, account1)
	require.NoError(t, err)
	require.Equal(t, account1.Balance, account1After.Balance)

	account2After, err := services.GetAccount(account2.ID)
	require.NotEmpty(t, account2)
	require.NoError(t, err)
	require.Equal(t, account2.Balance, account2After.Balance)

}

func createRandomUser(t *testing.T) models.User {
	createUserRequest := requests.CreateUserRequest{
		Username: util.RandomUsername(),
		Email:    util.RandomEmail(),
		FullName: util.RandomFullname(),
		Password: util.RandomPassword(),
	}

	createdTime := time.Now().Truncate(time.Nanosecond).Local()

	user, err := services.CreateUser(createUserRequest)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	require.Equal(t, createUserRequest.Username, user.Username)
	require.Equal(t, createUserRequest.Email, user.Email)
	require.NoError(t, util.CheckPassword(createUserRequest.Password, user.HashedPassword))
	require.Equal(t, createUserRequest.FullName, user.FullName)
	require.WithinDuration(t, createdTime, user.CreatedAt, time.Second)
	require.WithinDuration(t, createdTime, user.UpdatedAt, time.Second)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)

	res, err := services.GetUser(user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	require.Equal(t, user.Username, res.Username)
	require.Equal(t, user.FullName, res.FullName)
	require.Equal(t, user.Email, res.Email)
	require.Equal(t, user.HashedPassword, res.HashedPassword)
	require.WithinDuration(t, user.CreatedAt, res.CreatedAt, time.Millisecond)
	require.WithinDuration(t, user.UpdatedAt, res.UpdatedAt, time.Millisecond)
	require.Equal(t, user.DeletedAt, res.DeletedAt)
	require.True(t, res.DeletedAt.Time.IsZero())
}
