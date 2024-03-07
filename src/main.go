package main

import (
	"context"
	"crebito/config"
	"crebito/database"
	"crebito/server"
	"log"
)

func main() {
	cfg := config.New()

	if err := database.InitDatabase(cfg.DBHostname); err != nil {
		log.Panicf("error connecting database : error=%v", err)
	}
	defer database.DB.ExtractSession.Close(context.Background())
	defer database.DB.TransactionSession.Close(context.Background())
	defer database.DB.Driver.Close(context.Background())

	err := database.DB.Driver.VerifyConnectivity(context.Background())
	if err != nil {
		log.Println(err)
	}

	if err := server.Run(cfg.Mode, cfg.Port); err != nil {
		log.Panicf("error starting server : error=%v", err)
	}
}
