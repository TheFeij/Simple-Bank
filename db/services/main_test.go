package services

import (
	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var accountServices Services

func TestMain(m *testing.M) {
	dsn := "host=localhost user=root password=1234 dbname=simple_bank_test port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	accountServices = NewSQLServices(db)

	exitCode := m.Run()

	db.Exec("DELETE FROM accounts")
	db.Exec("DELETE FROM entries")
	db.Exec("DELETE FROM transfers")
	DB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(DB)

	os.Exit(exitCode)
}
