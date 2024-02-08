package services

import (
	"Simple-Bank/db/models"
	"Simple-Bank/requests"
	"Simple-Bank/responses"
	"database/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type SQLServices struct {
	DB *gorm.DB
}

func NewSQLServices(db *gorm.DB) Services {
	return &SQLServices{
		DB: db,
	}
}

func (services *SQLServices) CreateAccount(req requests.CreateAccountRequest) (responses.CreateAccountResponse, error) {
	newAccount := models.Accounts{
		Owner:   req.Owner,
		Balance: req.Balance,
	}

	if err := services.DB.Create(&newAccount).Error; err != nil {
		return responses.CreateAccountResponse{}, err
	}

	return responses.CreateAccountResponse{
		AccountID: newAccount.ID,
		CreatedAt: newAccount.CreatedAt,
		Owner:     newAccount.Owner,
		Balance:   newAccount.Balance,
	}, nil
}

func (services *SQLServices) DeleteAccount(id int64) (responses.GetAccountResponse, error) {
	var deletedAccount responses.GetAccountResponse

	if err := services.DB.
		Raw("SELECT id AS account_id, created_at, updated_at, deleted_at, owner, balance FROM accounts WHERE id = ?", id).
		Scan(&deletedAccount).Error; err != nil {
		return responses.GetAccountResponse{}, err
	}

	if err := services.DB.Delete(&models.Accounts{}, id).Error; err != nil {
		return responses.GetAccountResponse{}, err
	}

	return deletedAccount, nil
}

func (services *SQLServices) DepositMoney(req requests.DepositRequest) (responses.EntryResponse, error) {
	var newEntry models.Entries

	if err := services.DB.Transaction(func(tx *gorm.DB) error {
		newEntry = models.Entries{
			AccountID: req.AccountID,
			Amount:    req.Amount,
		}

		if err := tx.Create(&newEntry).Error; err != nil {
			return err
		}

		var account models.Accounts
		if err := tx.First(&account, req.AccountID).Error; err != nil {
			return nil
		}

		account.Balance += int64(req.Amount)
		if err := tx.Save(account).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return responses.EntryResponse{}, err
	}

	return responses.EntryResponse{
		EntryID:   newEntry.ID,
		AccountID: newEntry.AccountID,
		Amount:    newEntry.Amount,
		CreatedAt: newEntry.CreatedAt,
	}, nil
}

func (services *SQLServices) WithdrawMoney(req requests.WithdrawRequest) (responses.EntryResponse, error) {
	var newEntry models.Entries

	if err := services.DB.Transaction(func(tx *gorm.DB) error {
		newEntry = models.Entries{
			AccountID: req.AccountID,
			Amount:    -req.Amount,
		}

		if err := tx.Create(&newEntry).Error; err != nil {
			return err
		}

		var account models.Accounts
		if err := tx.First(&account, req.AccountID).Error; err != nil {
			return nil
		}

		account.Balance -= int64(req.Amount)
		if err := tx.Save(account).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return responses.EntryResponse{}, err
	}

	return responses.EntryResponse{
		EntryID:   newEntry.ID,
		AccountID: newEntry.AccountID,
		Amount:    newEntry.Amount,
		CreatedAt: newEntry.CreatedAt,
	}, nil
}

func (services *SQLServices) Transfer(req requests.TransferRequest) (responses.TransferResponse, error) {
	var newTransfer models.Transfers

	if err := services.DB.Transaction(func(tx *gorm.DB) error {
		var srcAccount, dstAccount models.Accounts

		// always acquire the lock of the account with the lower account id
		if req.FromAccountID < req.ToAccountID {
			var err error
			srcAccount, dstAccount, err = acquireLock(tx, req.FromAccountID, req.ToAccountID)
			if err != nil {
				return err
			}
		} else {
			var err error
			dstAccount, srcAccount, err = acquireLock(tx, req.ToAccountID, req.FromAccountID)
			if err != nil {
				return err
			}
		}

		srcAccount.Balance -= int64(req.Amount)
		dstAccount.Balance += int64(req.Amount)

		if err := tx.Save(&srcAccount).Error; err != nil {
			return err
		}
		if err := tx.Save(&dstAccount).Error; err != nil {
			return err
		}

		FromEntry := models.Entries{
			AccountID: req.FromAccountID,
			Amount:    -req.Amount,
		}
		ToEntry := models.Entries{
			AccountID: req.ToAccountID,
			Amount:    req.Amount,
		}
		if err := tx.Create(&FromEntry).Error; err != nil {
			return err
		}
		if err := tx.Create(&ToEntry).Error; err != nil {
			return err
		}

		newTransfer = models.Transfers{
			FromAccountID:   req.FromAccountID,
			ToAccountID:     req.ToAccountID,
			Amount:          req.Amount,
			OutgoingEntryID: FromEntry.ID,
			IncomingEntryID: ToEntry.ID,
		}

		if err := tx.Create(&newTransfer).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return responses.TransferResponse{}, err
	}

	return responses.TransferResponse{
		TransferID:   newTransfer.ID,
		SrcAccountID: newTransfer.FromAccountID,
		DstAccountID: newTransfer.ToAccountID,
		CreatedAt:    newTransfer.CreatedAt,
		Amount:       newTransfer.Amount,
	}, nil
}

func (services *SQLServices) ListAccounts(pageNumber int64, pageSize int8) (responses.ListAccountsResponse, error) {
	var accountsList responses.ListAccountsResponse

	offset := (pageNumber - 1) * int64(pageSize)
	res := services.DB.
		Raw("SELECT id AS account_id,"+
			" created_at,"+
			" deleted_at,"+
			" updated_at,"+
			" owner,"+
			" balance"+
			" FROM accounts LIMIT ? OFFSET ?", pageSize, offset).
		Scan(&accountsList.Accounts)

	if res.RowsAffected == 0 {
		return accountsList, sql.ErrNoRows
	}
	if res.Error != nil {
		return accountsList, res.Error
	}

	return accountsList, nil
}

func (services *SQLServices) GetAccount(id int64) (responses.GetAccountResponse, error) {
	var accountResponse responses.GetAccountResponse

	res := services.DB.
		Raw("SELECT id AS account_id, created_at, updated_at, deleted_at, owner, balance FROM accounts WHERE id = ?", id).
		Scan(&accountResponse)

	if res.RowsAffected == 0 {
		return accountResponse, sql.ErrNoRows
	}
	if res.Error != nil {
		return accountResponse, res.Error
	}

	return accountResponse, nil
}

func (services *SQLServices) GetTransfer(id int64) (responses.TransferResponse, error) {
	var transfer responses.TransferResponse

	if err := services.DB.
		Raw("SELECT id AS transfer_id,"+
			" to_account_id AS src_account_id,"+
			" from_account_id AS dst_account_id,"+
			" incoming_entry_id,"+
			" outgoing_entry_id,"+
			" created_at,"+
			" amount"+
			" FROM transfers WHERE id = ?", id).
		Scan(&transfer).Error; err != nil {
		return responses.TransferResponse{}, err
	}

	return transfer, nil
}

func (services *SQLServices) GetEntry(id int64) (responses.EntryResponse, error) {
	var entry responses.EntryResponse

	if err := services.DB.
		Raw("SELECT id AS entry_id, account_id, created_at, amount FROM entries WHERE id = ?", id).
		Scan(&entry).Error; err != nil {
		return responses.EntryResponse{}, err
	}

	return entry, nil
}

func (services *SQLServices) CreateUser(req requests.CreateUserRequest) (responses.UserInformationResponse, error) {
	newUser := models.User{
		Username:       req.Username,
		HashedPassword: "TODO",
		Email:          req.Email,
		FullName:       req.FullName,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		DeletedAt:      gorm.DeletedAt{},
	}

	if err := services.DB.Create(&newUser).Error; err != nil {
		return responses.UserInformationResponse{}, err
	}

	return responses.UserInformationResponse{
		Username:  newUser.Username,
		Email:     newUser.Email,
		FullName:  newUser.FullName,
		CreatedAt: newUser.CreatedAt.Local().Truncate(time.Second),
		UpdatedAt: newUser.UpdatedAt.Local().Truncate(time.Second),
		DeletedAt: newUser.DeletedAt.Time.Truncate(time.Second),
	}, nil

}

func (services *SQLServices) GetUser(username string) (responses.UserInformationResponse, error) {
	var user models.User
	var res responses.UserInformationResponse

	if err := services.DB.
		Raw("SELECT * FROM users WHERE username = ?", username).
		Scan(&user).Error; err != nil {
		return res, err
	}

	res = responses.UserInformationResponse{
		Username:  user.Username,
		Email:     user.Email,
		FullName:  user.FullName,
		CreatedAt: user.CreatedAt.Local().Truncate(time.Second),
		UpdatedAt: user.CreatedAt.Local().Truncate(time.Second),
	}

	if user.DeletedAt.Time.IsZero() {
		res.DeletedAt = user.DeletedAt.Time.Truncate(time.Second)
	} else {
		res.DeletedAt = user.DeletedAt.Time.Local().Truncate(time.Second)
	}

	return res, nil
}

func acquireLock(tx *gorm.DB, lowerAccountID, higherAccountID int64) (models.Accounts, models.Accounts, error) {
	var lowerAccount, higherAccount models.Accounts
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&lowerAccount, lowerAccountID).Error; err != nil {
		return lowerAccount, higherAccount, err
	}
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&higherAccount, higherAccountID).Error; err != nil {
		return lowerAccount, higherAccount, err
	}

	return lowerAccount, higherAccount, nil
}
