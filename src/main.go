package main

import (
	"crebito/config"
	"crebito/server"
	"log"
)

func main() {
	cfg := config.New()

	if err := server.Run(cfg.Mode, cfg.Port); err != nil {
		log.Panicf("error starting server : error=%v", err)
	}
}
