package main

import (
	"Simple-Bank/api"
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

	server := api.NewServer(db)
	if err := server.Start("localhost:8080"); err != nil {
		log.Fatalln("cannot start server: ", err)
	}
}
