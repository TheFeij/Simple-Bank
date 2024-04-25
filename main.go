package main

import (
	"Simple-Bank/api"
	"Simple-Bank/config"
	"Simple-Bank/db"
	"Simple-Bank/db/services"
	"Simple-Bank/token"
	"database/sql"
	"log"
)

func main() {
	configs, err := config.LoadConfig("./config", "config")
	if err != nil {
		log.Fatalln(err)
	}

	db.Init(configs.DatabaseSource)
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

	tokenMaker, err := token.NewPasetoMaker(configs.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("cannot create token maker: %w", err)
	}

	server, err := api.NewServer(&configs, services.NewSQLServices(db), tokenMaker)
	if err != nil {
		log.Fatalln("cannot create server: ", err)
	}

	if err = server.Start(configs.ServerHost + ":" + configs.ServerPort); err != nil {
		log.Fatalln("cannot start server: ", err)
	}
}
