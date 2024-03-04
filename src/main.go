package main

import (
	"context"
	"crebito/config"
	"crebito/database"
	"crebito/server"
	"log"
	"time"
)

func main() {
	cfg := config.New()

	if err := database.InitDatabase(cfg.URI, cfg.DBUsername, cfg.DBPassword); err != nil {
		log.Panicf("error connecting database : error=%v", err)
	}
	defer database.DBClient.CloseDB(context.Background())

	for {
		if err := database.DBClient.VerifyConnectivity(context.Background()); err != nil {
			log.Printf("error in database connectivity : error=%v\n", err)
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	if err := server.Run(cfg.Mode, cfg.Port); err != nil {
		log.Panicf("error starting server : error=%v", err)
	}
}
