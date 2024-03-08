package config

import "os"

type Config struct {
	Port       string
	DBHostname string
}

var cfg *Config

func init() {
	cfg = new(Config)

	cfg.Port = os.Getenv("PORT")
	cfg.DBHostname = os.Getenv("DB_HOSTNAME")
}

func New() *Config {
	return cfg
}
