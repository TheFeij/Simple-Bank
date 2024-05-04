package services

import (
	"Simple-Bank/db/models"
	"Simple-Bank/requests"
	"Simple-Bank/util"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
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

func (services *SQLServices) CreateAccount(owner string) (models.Account, error) {
	newAccount := models.Account{
		Owner:     owner,
		Balance:   0,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		DeletedAt: gorm.DeletedAt{},
	}

	if err := services.DB.Create(&newAccount).Error; err != nil {
		return newAccount, err
	}

	return newAccount, nil
}

func (services *SQLServices) DeleteAccount(id int64) (models.Account, error) {
	var deletedAccount models.Account

	if err := services.DB.Exec(
		"UPDATE accounts SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?", id).
		Error; err != nil {
		return deletedAccount, err
	}

	if err := services.DB.
		Raw("SELECT * FROM accounts WHERE id = ?", id).
		Scan(&deletedAccount).Error; err != nil {
		return deletedAccount, err
	}

	return deletedAccount, nil
}

func (services *SQLServices) DepositMoney(req requests.DepositRequest) (models.Entry, error) {
	var newEntry models.Entry

	if err := services.DB.Transaction(func(tx *gorm.DB) error {
		newEntry = models.Entry{
			AccountID: req.AccountID,
			Amount:    req.Amount,
		}

		if err := tx.Create(&newEntry).Error; err != nil {
			return err
		}

		var account models.Account
		if err := tx.First(&account, req.AccountID).Error; err != nil {
			return nil
		}

		account.Balance += int64(req.Amount)
		if err := tx.Save(account).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return newEntry, err
	}

	return newEntry, nil
}

func (services *SQLServices) WithdrawMoney(req requests.WithdrawRequest) (models.Entry, error) {
	var newEntry models.Entry

	if err := services.DB.Transaction(func(tx *gorm.DB) error {
		newEntry = models.Entry{
			AccountID: req.AccountID,
			Amount:    -req.Amount,
		}

		if err := tx.Create(&newEntry).Error; err != nil {
			return err
		}

		var account models.Account
		if err := tx.First(&account, req.AccountID).Error; err != nil {
			return nil
		}

		account.Balance -= int64(req.Amount)
		if err := tx.Save(account).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return newEntry, err
	}

	return newEntry, nil
}

func (services *SQLServices) Transfer(srcOwner string, req requests.TransferRequest) (models.Transfer, error) {
	var newTransfer models.Transfer

	if err := services.DB.Transaction(func(tx *gorm.DB) error {
		var srcAccount, dstAccount models.Account

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

		if srcAccount.Owner != srcOwner {
			err := fmt.Errorf("user is not the owner of the source account")
			return err
		}
		srcAccount.Balance -= int64(req.Amount)
		dstAccount.Balance += int64(req.Amount)

		if err := tx.Save(&srcAccount).Error; err != nil {
			return err
		}
		if err := tx.Save(&dstAccount).Error; err != nil {
			return err
		}

		FromEntry := models.Entry{
			AccountID: req.FromAccountID,
			Amount:    -req.Amount,
		}
		ToEntry := models.Entry{
			AccountID: req.ToAccountID,
			Amount:    req.Amount,
		}
		if err := tx.Create(&FromEntry).Error; err != nil {
			return err
		}
		if err := tx.Create(&ToEntry).Error; err != nil {
			return err
		}

		newTransfer = models.Transfer{
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
		return newTransfer, err
	}

	return newTransfer, nil
}

func (services *SQLServices) ListAccounts(owner string, pageNumber int64, pageSize int8) ([]models.Account, error) {
	var accountsList []models.Account

	offset := (pageNumber - 1) * int64(pageSize)
	res := services.DB.
		Raw("SELECT * FROM accounts WHERE owner = ? LIMIT ? OFFSET ?", owner, pageSize, offset).
		Scan(&accountsList)

	if res.RowsAffected == 0 {
		return accountsList, sql.ErrNoRows
	}
	if res.Error != nil {
		return accountsList, res.Error
	}

	return accountsList, nil
}

func (services *SQLServices) GetAccount(id int64) (models.Account, error) {
	var account models.Account

	res := services.DB.
		Raw("SELECT * FROM accounts WHERE id = ?", id).
		Scan(&account)

	if res.RowsAffected == 0 {
		return account, sql.ErrNoRows
	}
	if res.Error != nil {
		return account, res.Error
	}

	return account, nil
}

func (services *SQLServices) GetTransfer(id int64) (models.Transfer, error) {
	var transfer models.Transfer

	if err := services.DB.First(&transfer, id).Error; err != nil {
		return transfer, err
	}

	return transfer, nil
}

func (services *SQLServices) GetEntry(id int64) (models.Entry, error) {
	var entry models.Entry

	if err := services.DB.First(&entry, id).Error; err != nil {
		return entry, err
	}

	return entry, nil
}

func (services *SQLServices) CreateUser(req requests.CreateUserRequest) (models.User, error) {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return models.User{}, err
	}

	newUser := models.User{
		Username:       req.Username,
		Email:          req.Email,
		FullName:       req.FullName,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
		DeletedAt:      gorm.DeletedAt{},
	}

	if err := services.DB.Create(&newUser).Error; err != nil {
		return models.User{}, err
	}

	return newUser, nil
}

func (services *SQLServices) GetUser(username string) (models.User, error) {
	var user models.User

	if err := services.DB.Model(&models.User{}).Where("username = ?", username).First(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (services *SQLServices) CreateSession(session models.Session) (models.Session, error) {
	if err := services.DB.Create(&session).Error; err != nil {
		return models.Session{}, err
	}

	return session, nil
}

func (services *SQLServices) GetSession(id uuid.UUID) (models.Session, error) {
	var session models.Session

	if err := services.DB.
		Raw("SELECT * FROM sessions WHERE id = ?", id).
		Scan(&session).Error; err != nil {
		return models.Session{}, err
	}

	return session, nil
}

// UpdateUser updates the user information in the database based on the provided request.
// It hashes the password if provided, and updates the fullname and email fields if they are not nil.
// It returns the updated user model and any error encountered.
func (services *SQLServices) UpdateUser(req UpdateUserRequest) (models.User, error) {
	updateData := map[string]interface{}{}

	if req.Password != nil {
		hashedPassword, err := util.HashPassword(*req.Password)
		if err != nil {
			return models.User{}, err
		}
		updateData["password"] = hashedPassword
	}
	if req.Fullname != nil {
		updateData["fullname"] = req.Fullname
	}
	if req.Email != nil {
		updateData["email"] = req.Email
	}

	var user models.User
	if err := services.DB.
		Model(&models.User{}).
		Where("username = ?", req.Username).
		Updates(updateData).
		First(&user).Error; err != nil {
		return models.User{}, err
	}

	return user, nil
}

func acquireLock(tx *gorm.DB, lowerAccountID, higherAccountID int64) (models.Account, models.Account, error) {
	var lowerAccount, higherAccount models.Account
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
