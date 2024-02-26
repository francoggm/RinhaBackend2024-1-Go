package main

import (
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
	defer database.DBClient.CloseDB(database.DBContext)

	for {
		if err := database.DBClient.VerifyConnectivity(database.DBContext); err != nil {
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
