package seed

import (
	"Simple-Bank/db/models"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	account1 := models.Accounts{
		Owner:   "Abolfazl Moradi Feijani",
		Balance: 3000,
	}
	account2 := models.Accounts{
		Owner:   "Ali Mohammadi",
		Balance: 1000,
	}
	account3 := models.Accounts{
		Owner:   "Ahmad Babaee",
		Balance: 5000,
	}

	db.Create(&account1)
	db.Create(&account2)
	db.Create(&account3)

	entry1 := models.Entry{
		AccountID: 3,
		Amount:    500,
	}
	entry2 := models.Entry{
		AccountID: 1,
		Amount:    -200,
	}
	entry3 := models.Entry{
		AccountID: 3,
		Amount:    -300,
	}
	entry4 := models.Entry{
		AccountID: 3,
		Amount:    200,
	}
	entry5 := models.Entry{
		AccountID: 2,
		Amount:    300,
	}
	entry6 := models.Entry{
		AccountID: 2,
		Amount:    -500,
	}

	db.Create(&entry1)
	db.Create(&entry2)
	db.Create(&entry3)
	db.Create(&entry4)
	db.Create(&entry5)
	db.Create(&entry6)

	transfer1 := models.Transfers{
		FromAccountID: 2,
		ToAccountID:   3,
		Amount:        500,
	}
	transfer2 := models.Transfers{
		FromAccountID: 1,
		ToAccountID:   3,
		Amount:        200,
	}
	transfer3 := models.Transfers{
		FromAccountID: 3,
		ToAccountID:   2,
		Amount:        300,
	}

	db.Create(&transfer1)
	db.Create(&transfer2)
	db.Create(&transfer3)
}
