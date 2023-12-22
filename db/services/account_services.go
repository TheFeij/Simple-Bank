package services

import (
	"Simple-Bank/db/models"
	"Simple-Bank/requests"
	"Simple-Bank/responses"
	"gorm.io/gorm"
)

type AccountServices struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *AccountServices {
	return &AccountServices{
		DB: db,
	}
}

func (accountServices *AccountServices) CreateAccount(req requests.CreateAccountRequest) (responses.CreateAccountResponse, error) {
	newAccount := models.Accounts{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	if err := accountServices.DB.Create(&newAccount).Error; err != nil {
		return responses.CreateAccountResponse{}, err
	}

	response := responses.CreateAccountResponse{
		AccountID: uint64(newAccount.ID),
		Owner:     newAccount.Owner,
		Balance:   newAccount.Balance,
		CreatedAt: newAccount.CreatedAt,
		Currency:  newAccount.Currency,
	}
	return response, nil
}

func (accountServices *AccountServices) DeleteAccount(id uint64) (responses.GetAccountResponse, error) {
	var account models.Accounts

	if err := accountServices.DB.First(&account, id).Error; err != nil {
		return responses.GetAccountResponse{}, err
	}

	if err := accountServices.DB.Delete(&account).Error; err != nil {
		return responses.GetAccountResponse{}, err
	}

	return responses.GetAccountResponse{
		AccountID: uint64(account.ID),
		Owner:     account.Owner,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt,
		Currency:  account.Currency,
	}, nil
}

func (accountServices *AccountServices) DepositMoney(req requests.DepositRequest) (models.Entries, error) {
	var newEntry models.Entries

	if err := accountServices.DB.Transaction(func(tx *gorm.DB) error {
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
		return models.Entries{}, err
	}

	return newEntry, nil
}

func (accountServices *AccountServices) WithdrawMoney(req requests.WithdrawRequest) (models.Entries, error) {
	var newEntry models.Entries

	if err := accountServices.DB.Transaction(func(tx *gorm.DB) error {
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
		return models.Entries{}, err
	}

	return newEntry, nil
}

/*
should validate data prior to db call but here i ignore them for less complexity
*/
func (accountServices *AccountServices) Transfer(req requests.TransferRequest) (models.Transfers, error) {
	var newTransfer models.Transfers

	if err := accountServices.DB.Transaction(func(tx *gorm.DB) error {
		err := tx.Exec(`set transaction isolation level repeatable read`).Error
		if err != nil {
			return err
		}
		newTransfer = models.Transfers{
			FromAccountID: req.FromAccountID,
			ToAccountID:   req.ToAccountID,
			Amount:        req.Amount,
		}

		if err := tx.Create(&newTransfer).Error; err != nil {
			return err
		}

		var srcAccount, dstAccount models.Accounts
		if err := tx.First(&srcAccount, req.FromAccountID).Error; err != nil {
			return err
		}
		if err := tx.First(&dstAccount, req.ToAccountID).Error; err != nil {
			return err
		}

		srcAccount.Balance -= uint64(req.Amount)
		dstAccount.Balance += uint64(req.Amount)

		if err := tx.Save(&srcAccount).Error; err != nil {
			return err
		}
		if err := tx.Save(&dstAccount).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return models.Transfers{}, err
	}

	return newTransfer, nil
}

func (accountServices *AccountServices) ListAccounts(limit int) ([]responses.ListAccountsResponse, error) {
	var accounts []models.Accounts

	if err := accountServices.DB.First(&accounts).Limit(limit).Error; err != nil {
		return []responses.ListAccountsResponse{}, err
	}

	var accountsList []responses.ListAccountsResponse
	for _, account := range accounts {
		response := responses.ListAccountsResponse{
			Owner:     account.Owner,
			Currency:  account.Currency,
			Balance:   account.Balance,
			AccountID: uint64(account.ID),
			CreatedAt: account.CreatedAt,
			UpdateAt:  account.UpdatedAt,
		}

		accountsList = append(accountsList, response)
	}

	return accountsList, nil
}

func (accountServices *AccountServices) GetAccount(id uint64) (responses.GetAccountResponse, error) {
	var account models.Accounts

	if err := accountServices.DB.Find(&account, id).Error; err != nil {
		return responses.GetAccountResponse{}, err
	}

	response := responses.GetAccountResponse{
		Owner:     account.Owner,
		Currency:  account.Currency,
		Balance:   account.Balance,
		AccountID: uint64(account.ID),
		CreatedAt: account.CreatedAt,
		UpdatedAt: account.UpdatedAt,
	}

	return response, nil
}
