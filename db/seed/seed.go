package seed

import (
	"Simple-Bank/db/models"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	account1 := models.Accounts{
		Owner:    "Abolfazl Moradi Feijani",
		Balance:  3000,
		Currency: "USD",
	}
	account2 := models.Accounts{
		Owner:    "Ali Mohammadi",
		Balance:  1000,
		Currency: "IRL",
	}
	account3 := models.Accounts{
		Owner:    "Ahmad Babaee",
		Balance:  5000,
		Currency: "USD",
	}

	db.Create(&account1)
	db.Create(&account2)
	db.Create(&account3)

	entry1 := models.Entries{
		AccountID: 3,
		Amount:    500,
	}
	entry2 := models.Entries{
		AccountID: 1,
		Amount:    -200,
	}
	entry3 := models.Entries{
		AccountID: 3,
		Amount:    -300,
	}

	db.Create(&entry1)
	db.Create(&entry2)
	db.Create(&entry3)

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
