package services

import (
	"Simple-Bank/db/models"
	"Simple-Bank/requests"
	"Simple-Bank/responses"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Services struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Services {
	return &Services{
		DB: db,
	}
}

func (services *Services) CreateAccount(req requests.CreateAccountRequest) (responses.CreateAccountResponse, error) {
	newAccount := models.Accounts{
		Owner:   req.Owner,
		Balance: req.Balance,
	}

	if err := services.DB.Create(&newAccount).Error; err != nil {
		return responses.CreateAccountResponse{}, err
	}

	return responses.CreateAccountResponse{
		AccountID: uint64(newAccount.ID),
		CreatedAt: newAccount.CreatedAt,
		Owner:     newAccount.Owner,
		Balance:   newAccount.Balance,
	}, nil
}

func (services *Services) DeleteAccount(id uint64) (responses.GetAccountResponse, error) {
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

func (services *Services) DepositMoney(req requests.DepositRequest) (responses.EntryResponse, error) {
	var newEntry models.Entries

	if err := services.DB.Transaction(func(tx *gorm.DB) error {
		newEntry = models.Entries{
			AccountID: req.AccountID,
			Amount:    int64(req.Amount),
		}

		if err := tx.Create(&newEntry).Error; err != nil {
			return err
		}

		var account models.Accounts
		if err := tx.First(&account, req.AccountID).Error; err != nil {
			return nil
		}

		account.Balance += uint64(req.Amount)
		if err := tx.Save(account).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return responses.EntryResponse{}, err
	}

	return responses.EntryResponse{
		EntryID:   uint64(newEntry.ID),
		AccountID: newEntry.AccountID,
		Amount:    newEntry.Amount,
		CreatedAt: newEntry.CreatedAt,
	}, nil
}

func (services *Services) WithdrawMoney(req requests.WithdrawRequest) (responses.EntryResponse, error) {
	var newEntry models.Entries

	if err := services.DB.Transaction(func(tx *gorm.DB) error {
		newEntry = models.Entries{
			AccountID: req.AccountID,
			Amount:    -1 * int64(req.Amount),
		}

		if err := tx.Create(&newEntry).Error; err != nil {
			return err
		}

		var account models.Accounts
		if err := tx.First(&account, req.AccountID).Error; err != nil {
			return nil
		}

		account.Balance -= uint64(req.Amount)
		if err := tx.Save(account).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return responses.EntryResponse{}, err
	}

	return responses.EntryResponse{
		EntryID:   uint64(newEntry.ID),
		AccountID: newEntry.AccountID,
		Amount:    newEntry.Amount,
		CreatedAt: newEntry.CreatedAt,
	}, nil
}

func (services *Services) Transfer(req requests.TransferRequest) (responses.TransferResponse, error) {
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

		srcAccount.Balance -= uint64(req.Amount)
		dstAccount.Balance += uint64(req.Amount)

		if err := tx.Save(&srcAccount).Error; err != nil {
			return err
		}
		if err := tx.Save(&dstAccount).Error; err != nil {
			return err
		}

		FromEntry := models.Entries{
			AccountID: req.FromAccountID,
			Amount:    -1 * int64(req.Amount),
		}
		ToEntry := models.Entries{
			AccountID: req.ToAccountID,
			Amount:    int64(req.Amount),
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
			OutgoingEntryID: uint64(FromEntry.ID),
			IncomingEntryID: uint64(ToEntry.ID),
		}

		if err := tx.Create(&newTransfer).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return responses.TransferResponse{}, err
	}

	return responses.TransferResponse{
		TransferID:   uint64(newTransfer.ID),
		SrcAccountID: newTransfer.FromAccountID,
		DstAccountID: newTransfer.ToAccountID,
		CreatedAt:    newTransfer.CreatedAt,
		Amount:       uint32(int64(newTransfer.Amount)),
	}, nil
}

func (services *Services) ListAccounts(limit int) (responses.ListAccountsResponse, error) {
	var accountsList responses.ListAccountsResponse

	if err := services.DB.
		Raw("SELECT id AS account_id, created_at, deleted_at, updated_at, owner, balance FROM accounts LIMIT ?", limit).
		Scan(&accountsList.Accounts).Error; err != nil {
		return responses.ListAccountsResponse{}, err
	}

	return accountsList, nil
}

func (services *Services) GetAccount(id uint64) (responses.GetAccountResponse, error) {
	var accountResponse responses.GetAccountResponse

	if err := services.DB.
		Raw("SELECT id AS account_id, created_at, updated_at, deleted_at, owner, balance FROM accounts WHERE id = ?", id).
		Scan(&accountResponse).Error; err != nil {
		return responses.GetAccountResponse{}, err
	}

	return accountResponse, nil
}

func (services *Services) GetTransfer(id uint64) (responses.TransferResponse, error) {
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

func (services *Services) GetEntry(id uint64) (responses.EntryResponse, error) {
	var entry responses.EntryResponse

	if err := services.DB.
		Raw("SELECT id AS entry_id, account_id, created_at, amount FROM entries WHERE id = ?", id).
		Scan(&entry).Error; err != nil {
		return responses.EntryResponse{}, err
	}

	return entry, nil
}

func acquireLock(tx *gorm.DB, lowerAccountID, higherAccountID uint64) (models.Accounts, models.Accounts, error) {
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
