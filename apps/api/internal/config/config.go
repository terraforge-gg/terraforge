package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                  string
	HostPort             string
	FrontendUrl          string
	AuthUrl              string
	DatabaseUrl          string
	MeiliSearchHostUrl   string
	MeiliSearchMasterKey string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Env:                  os.Getenv("ENV"),
		HostPort:             os.Getenv("HOST_PORT"),
		FrontendUrl:          os.Getenv("FRONTEND_URL"),
		AuthUrl:              os.Getenv("AUTH_URL"),
		DatabaseUrl:          os.Getenv("DATABASE_URL"),
		MeiliSearchHostUrl:   os.Getenv("MEILISEARCH_HOST_URL"),
		MeiliSearchMasterKey: os.Getenv("MEILISEARCH_MASTER_KEY"),
	}
}
