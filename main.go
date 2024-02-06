package main

import (
	"Simple-Bank/api"
	"Simple-Bank/config"
	"Simple-Bank/db"
	"database/sql"
	"log"
)

func main() {
	configs, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatalln(err)
	}

	db.Init(configs.Database.Source)
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
	if err := server.Start(configs.Server.Host + ":" + configs.Server.Port); err != nil {
		log.Fatalln("cannot start server: ", err)
	}
}
