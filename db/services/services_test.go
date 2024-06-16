package services

import (
	"Simple-Bank/db/models"
	"Simple-Bank/requests"
	"Simple-Bank/util"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
	"testing"
	"time"
)

func createAccount(t *testing.T, owner string) models.Account {

	createdTime := time.Now().Truncate(time.Nanosecond).Local()

	account, err := services.CreateAccount(owner)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, owner, account.Owner)
	require.Equal(t, int64(0), account.Balance)
	require.True(t, account.ID > 0)
	require.WithinDuration(t, createdTime, account.CreatedAt, time.Second)
	require.WithinDuration(t, createdTime, account.UpdatedAt, time.Second)
	require.True(t, account.DeletedAt.Time.IsZero())

	return account
}

func TestCreateAccount(t *testing.T) {
	user := createRandomUser(t)
	createAccount(t, user.Username)
}

func TestGetAccount(t *testing.T) {
	t.Run("UserFound", func(t *testing.T) {
		user := createRandomUser(t)
		account := createAccount(t, user.Username)

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
	})
	t.Run("UserNotFound", func(t *testing.T) {
		response, err := services.GetAccount(util.RandomID())
		require.Error(t, err)
		require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		require.Empty(t, response)
	})
}

func TestDeleteAccount(t *testing.T) {
	t.Run("AccountDeletedSuccessfully", func(t *testing.T) {
		user := createRandomUser(t)
		account := createAccount(t, user.Username)

		response, err := services.DeleteAccount(account.ID)
		require.NoError(t, err)
		require.NotEmpty(t, response)
		require.Equal(t, account.ID, response.ID)
		require.Equal(t, account.Balance, response.Balance)
		require.Equal(t, account.Owner, response.Owner)
		require.WithinDuration(t, account.CreatedAt, response.CreatedAt, time.Second)
		require.WithinDuration(t, account.UpdatedAt, response.UpdatedAt, time.Second)
		require.WithinDuration(t, time.Now(), response.DeletedAt.Time, time.Second)
		require.True(t, response.DeletedAt.Valid)

		response, err = services.GetAccount(account.ID)
		require.Error(t, err)
		require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		require.Empty(t, response)
	})
	t.Run("AccountNotFound", func(t *testing.T) {
		response, err := services.DeleteAccount(util.RandomID())
		require.Error(t, err)
		require.ErrorIs(t, err, ErrSrcAccountNotFound)
		require.Empty(t, response)
	})
}

func TestGetAccountsList(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		user := createRandomUser(t)

		createdAccounts := make([]models.Account, 5)
		for i := 0; i < 5; i++ {
			createdAccounts[i] = createAccount(t, user.Username)
		}

		accounts, err := services.ListAccounts(ListAccountsRequest{
			Owner:      user.Username,
			PageSize:   1,
			PageNumber: 5,
		})

		require.NoError(t, err)
		require.NotEmpty(t, accounts)
		require.Len(t, accounts, 5)

		for i, account := range accounts {
			require.NotEmpty(t, account)

			require.Equal(t, user.Username, account.Owner)
			require.Equal(t, createdAccounts[i].Balance, account.Balance)
			require.True(t, account.ID > 0)
			require.WithinDuration(t, createdAccounts[i].CreatedAt, account.CreatedAt, time.Second)
			require.WithinDuration(t, createdAccounts[i].UpdatedAt, account.UpdatedAt, time.Second)
			require.True(t, account.DeletedAt.Time.IsZero())
		}
	})
	t.Run("NoAccountsFound", func(t *testing.T) {
		accounts, err := services.ListAccounts(ListAccountsRequest{
			Owner:      util.RandomUsername(),
			PageSize:   1,
			PageNumber: 5,
		})
		require.Empty(t, accounts)
		require.Error(t, err)
		require.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	})
}

func TestTransfer(t *testing.T) {
	t.Run("Concurrent Transfers From 1 to 2", func(t *testing.T) {
		user1 := createRandomUser(t)
		user2 := createRandomUser(t)
		srcOwner := user1.Username
		account1 := createAccount(t, user1.Username)
		account2 := createAccount(t, user2.Username)

		concurrentTransactions := 20
		var amount int32 = 10

		srcInitialBalance := int32(concurrentTransactions) * amount

		// first deposit the amount of money to the src account before transferring
		_, err := services.DepositMoney(DepositRequest{
			Owner:     user1.Username,
			AccountID: account1.ID,
			Amount:    srcInitialBalance,
		})
		require.NoError(t, err)

		errorsChan := make(chan error)
		resultsChan := make(chan models.Transfer)

		for i := 0; i < concurrentTransactions; i++ {
			go func(chan models.Transfer, chan error) {
				transferRequest := TransferRequest{
					Owner:         srcOwner,
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

			require.True(t, int64((i+1)*int(amount)) <= toAccount.Balance)
			require.True(t, int64(int(srcInitialBalance)-(i+1)*int(amount)) >= fromAccount.Balance)
		}
	})
	t.Run("Concurrent Transfers From 1 to 2 and Reverse", func(t *testing.T) {
		user1 := createRandomUser(t)
		user2 := createRandomUser(t)
		account1 := createAccount(t, user1.Username)
		account2 := createAccount(t, user2.Username)

		concurrentTransactions := 20
		var amount int32 = 10

		initialBalance := int32(concurrentTransactions/2) * amount

		// first deposit the amount of money to the src and dst account before transferring
		_, err := services.DepositMoney(DepositRequest{
			Owner:     user1.Username,
			AccountID: account1.ID,
			Amount:    initialBalance,
		})
		require.NoError(t, err)

		_, err = services.DepositMoney(DepositRequest{
			Owner:     user2.Username,
			AccountID: account2.ID,
			Amount:    initialBalance,
		})
		require.NoError(t, err)

		errorsChan := make(chan error)
		resultsChan := make(chan models.Transfer)

		for i := 0; i < concurrentTransactions; i++ {
			reverse := i%2 == 0
			go func(chan models.Transfer, chan error, bool) {
				srcOwner := account1.Owner
				fromAccountID, toAccountID := account1.ID, account2.ID
				if reverse {
					fromAccountID, toAccountID = toAccountID, fromAccountID
					srcOwner = account2.Owner
				}
				transferRequest := TransferRequest{
					Owner:         srcOwner,
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
		require.Equal(t, int64(initialBalance), account1After.Balance)

		account2After, err := services.GetAccount(account2.ID)
		require.NotEmpty(t, account2)
		require.NoError(t, err)
		require.Equal(t, int64(initialBalance), account2After.Balance)
	})
	t.Run("Source Account Not Found", func(t *testing.T) {
		user1 := createRandomUser(t)
		user2 := createRandomUser(t)

		account2 := createAccount(t, user2.Username)

		req := TransferRequest{
			Owner:         user1.Username,
			FromAccountID: util.RandomID(),
			ToAccountID:   account2.ID,
			Amount:        200,
		}

		res, err := services.Transfer(req)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrSrcAccountNotFound)
		require.Empty(t, res)
	})
	t.Run("Destination Account Not Found", func(t *testing.T) {
		user1 := createRandomUser(t)

		account1 := createAccount(t, user1.Username)

		req := TransferRequest{
			Owner:         user1.Username,
			FromAccountID: account1.ID,
			ToAccountID:   util.RandomID(),
			Amount:        200,
		}

		res, err := services.Transfer(req)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrDstAccountNotFound)
		require.Empty(t, res)
	})
	t.Run("UnAuthorized Owner", func(t *testing.T) {
		user1 := createRandomUser(t)
		user2 := createRandomUser(t)

		account2 := createAccount(t, user2.Username)
		account1 := createAccount(t, user1.Username)

		req := TransferRequest{
			Owner:         util.RandomUsername(),
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Amount:        200,
		}

		res, err := services.Transfer(req)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrUnAuthorizedTransfer)
		require.Empty(t, res)
	})
	t.Run("Not Enough Money", func(t *testing.T) {
		user1 := createRandomUser(t)
		user2 := createRandomUser(t)

		account2 := createAccount(t, user2.Username)
		account1 := createAccount(t, user1.Username)

		req := TransferRequest{
			Owner:         user1.Username,
			FromAccountID: account1.ID,
			ToAccountID:   account2.ID,
			Amount:        200,
		}

		res, err := services.Transfer(req)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNotEnoughMoney)
		require.Empty(t, res)
	})
}

func TestDepositMoney(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		user := createRandomUser(t)
		account := createAccount(t, user.Username)

		req := DepositRequest{
			Owner:     user.Username,
			AccountID: account.ID,
			Amount:    200,
		}

		start := time.Now()

		res, err := services.DepositMoney(req)
		require.NoError(t, err)
		require.NotEmpty(t, res)

		require.Equal(t, req.AccountID, res.AccountID)
		require.Equal(t, req.Amount, res.Amount)
		require.WithinDuration(t, start, res.CreatedAt, time.Second)
		require.Zero(t, res.DeletedAt)

		account, err = services.GetAccount(req.AccountID)
		require.NoError(t, err)
		require.NotEmpty(t, account)

		require.Equal(t, int64(req.Amount), account.Balance)
	})
	t.Run("Invalid Owner", func(t *testing.T) {
		user := createRandomUser(t)
		account := createAccount(t, user.Username)

		req := DepositRequest{
			Owner:     "invalid owner",
			AccountID: account.ID,
			Amount:    200,
		}

		res, err := services.DepositMoney(req)
		require.Error(t, err)
		require.ErrorIs(t, ErrUnAuthorizedDeposit, err)
		require.Empty(t, res)
	})
	t.Run("Account Not Found", func(t *testing.T) {
		user := createRandomUser(t)

		req := DepositRequest{
			Owner:     user.Username,
			AccountID: util.RandomID(),
			Amount:    200,
		}

		res, err := services.DepositMoney(req)
		require.Error(t, err)
		require.ErrorIs(t, ErrSrcAccountNotFound, err)
		require.Empty(t, res)
	})
}

func TestWithdrawMoney(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		user := createRandomUser(t)
		account := createAccount(t, user.Username)

		res, err := services.DepositMoney(DepositRequest{
			Owner:     user.Username,
			AccountID: account.ID,
			Amount:    200,
		})
		require.NoError(t, err)

		req := WithdrawRequest{
			Owner:     user.Username,
			AccountID: account.ID,
			Amount:    200,
		}

		start := time.Now()

		res, err = services.WithdrawMoney(req)
		require.NoError(t, err)
		require.NotEmpty(t, res)

		require.Equal(t, req.AccountID, res.AccountID)
		require.Equal(t, -req.Amount, res.Amount)
		require.WithinDuration(t, start, res.CreatedAt, time.Second)
		require.Zero(t, res.DeletedAt)

		account, err = services.GetAccount(req.AccountID)
		require.NoError(t, err)
		require.NotEmpty(t, account)

		require.Equal(t, int64(0), account.Balance)
	})
	t.Run("Invalid Owner", func(t *testing.T) {
		user := createRandomUser(t)
		account := createAccount(t, user.Username)

		req := WithdrawRequest{
			Owner:     "invalid owner",
			AccountID: account.ID,
			Amount:    200,
		}

		res, err := services.WithdrawMoney(req)
		require.Error(t, err)
		require.ErrorIs(t, ErrUnAuthorizedWithdraw, err)
		require.Empty(t, res)
	})
	t.Run("Account Not Found", func(t *testing.T) {
		user := createRandomUser(t)

		req := WithdrawRequest{
			Owner:     user.Username,
			AccountID: util.RandomID(),
			Amount:    200,
		}

		res, err := services.WithdrawMoney(req)
		require.Error(t, err)
		require.ErrorIs(t, ErrSrcAccountNotFound, err)
		require.Empty(t, res)
	})
	t.Run("Not Enough Money", func(t *testing.T) {
		user := createRandomUser(t)
		account := createAccount(t, user.Username)

		req := WithdrawRequest{
			Owner:     user.Username,
			AccountID: account.ID,
			Amount:    200,
		}

		res, err := services.WithdrawMoney(req)
		require.Error(t, err)
		require.ErrorIs(t, ErrNotEnoughMoney, err)
		require.Empty(t, res)
	})
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
	var user models.User
	t.Run("UserCreated", func(t *testing.T) {
		user = createRandomUser(t)
	})
	t.Run("DuplicateUsername", func(t *testing.T) {
		user, err := services.CreateUser(requests.CreateUserRequest{
			Username: user.Username,
			Email:    util.RandomEmail(),
			FullName: util.RandomFullname(),
			Password: util.RandomPassword(),
		})
		require.Error(t, err)
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		require.Equal(t, "users_pkey", pgErr.ConstraintName)
		require.Empty(t, user)
	})
	t.Run("DuplicateEmail", func(t *testing.T) {
		user, err := services.CreateUser(requests.CreateUserRequest{
			Username: util.RandomUsername(),
			Email:    user.Email,
			FullName: util.RandomFullname(),
			Password: util.RandomPassword(),
		})
		require.Error(t, err)
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		require.Equal(t, "users_email_key", pgErr.ConstraintName)
		require.Empty(t, user)
	})
}

func TestGetUser(t *testing.T) {
	user := createRandomUser(t)

	t.Run("UserFound", func(t *testing.T) {
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
	})
	t.Run("UserNotFound", func(t *testing.T) {
		res, err := services.GetUser(util.RandomUsername())
		require.Error(t, err)
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
		require.Empty(t, res)
	})
}

func createSession(t *testing.T) models.Session {
	user := createRandomUser(t)

	session := models.Session{
		ID:           uuid.New(),
		Username:     user.Username,
		RefreshToken: "refresh token",
		UserAgent:    "user agent",
		ClientIP:     util.RandomIP(),
		IsBlocked:    false,
		CreatedAt:    time.Now(),
		ExpiresAt:    time.Now(),
		DeletedAt:    gorm.DeletedAt{},
	}

	returnedSession, err := services.CreateSession(session)
	require.NoError(t, err)

	require.Equal(t, session.ID, returnedSession.ID)
	require.Equal(t, session.UserAgent, returnedSession.UserAgent)
	require.Equal(t, session.ClientIP, returnedSession.ClientIP)
	require.Equal(t, session.Username, returnedSession.Username)
	require.WithinDuration(t, session.CreatedAt, returnedSession.CreatedAt, time.Second)
	require.WithinDuration(t, session.ExpiresAt, returnedSession.ExpiresAt, time.Second)
	require.Equal(t, session.DeletedAt, returnedSession.DeletedAt)
	require.Equal(t, session.IsBlocked, returnedSession.IsBlocked)

	return session
}

func TestSQLServices_CreateSession(t *testing.T) {
	createSession(t)
}

func TestSQLServices_Session(t *testing.T) {
	session := createSession(t)

	returnedSession, err := services.GetSession(session.ID)
	require.NoError(t, err)

	require.Equal(t, session.ID, returnedSession.ID)
	require.Equal(t, session.UserAgent, returnedSession.UserAgent)
	require.Equal(t, session.ClientIP, returnedSession.ClientIP)
	require.Equal(t, session.Username, returnedSession.Username)
	require.WithinDuration(t, session.CreatedAt, returnedSession.CreatedAt, time.Second)
	require.WithinDuration(t, session.ExpiresAt, returnedSession.ExpiresAt, time.Second)
	require.Equal(t, session.DeletedAt, returnedSession.DeletedAt)
	require.Equal(t, session.IsBlocked, returnedSession.IsBlocked)
}
