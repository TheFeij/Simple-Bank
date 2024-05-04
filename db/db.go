package db

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init(source string) error {
	db, err := gorm.Open(postgres.Open(source), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("cannot initialize dataabse: %w", err)
	}

	DB = db

	return nil
}

func GetDB() *gorm.DB {
	return DB
}
