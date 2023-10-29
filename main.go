package main

import (
	"Simple-Bank/db"
	"database/sql"
	"log"
)

func main() {
	db.Init()
	db := db.GetDB()
	DB, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	defer func(DB *sql.DB) {
		err := DB.Close()
		if err != nil {
			log.Fatalln(err)
		}
	}(DB)
}
