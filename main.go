package main

import (
	"Simple-Bank/api"
	"Simple-Bank/config"
	"Simple-Bank/db"
	"Simple-Bank/db/services"
	"Simple-Bank/token"
	"database/sql"
	"gorm.io/gorm"
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
		log.Fatalf("cannot create token maker: %v", err)
	}

	runGinServer(configs, tokenMaker, db)
}

func runGinServer(config config.Config, tokenMaker token.Maker, db *gorm.DB) {
	server, err := api.NewServer(&config, services.NewSQLServices(db), tokenMaker)
	if err != nil {
		log.Fatalln("cannot create server: ", err)
	}

	if err = server.Start(config.HTTPServerHost + ":" + config.HTTPServerPort); err != nil {
		log.Fatalln("cannot start server: ", err)
	}
}
