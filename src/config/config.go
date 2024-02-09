package config

import "os"

type Config struct {
	Mode       string
	Port       string
	DBUsername string
	DBPassword string
}

var cfg *Config

func init() {
	cfg = new(Config)

	cfg.Mode = os.Getenv("MODE")
	cfg.Port = os.Getenv("PORT")
	cfg.DBUsername = os.Getenv("NEO4J_USERNAME")
	cfg.DBPassword = os.Getenv("NEO4J_PASSWORD")
}

func New() *Config {
	return cfg
}
