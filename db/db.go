package db

import (
	models2 "Simple-Bank/db/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Init() {
	dsn := "host=localhost user=root password=david1380 dbname=simple_bank port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	db.AutoMigrate(&models2.Accounts{})
	db.AutoMigrate(&models2.Entries{})
	db.AutoMigrate(&models2.Transfers{})

	DB = db
}

func GetDB() *gorm.DB {
	return DB
}
