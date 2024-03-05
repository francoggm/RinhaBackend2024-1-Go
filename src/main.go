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
	defer database.DBClient.CloseDB(context.Background())

	if err := server.Run(cfg.Mode, cfg.Port); err != nil {
		log.Panicf("error starting server : error=%v", err)
	}
}
