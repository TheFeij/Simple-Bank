package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Init(source string) {
	db, err := gorm.Open(postgres.Open(source), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	DB = db
}

func GetDB() *gorm.DB {
	return DB
}
