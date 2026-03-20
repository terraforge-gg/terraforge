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
	SeedDb               bool
	CfRegion             string
	CfAccountId          string
	CfAccessKeyId        string
	CfSecretAccessKey    string
	CdnUrl               string
	R2Bucket             string
	MeiliSearchHostUrl   string
	MeiliSearchMasterKey string
	RedisUrl             string
	RedisPassword        string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		Env:                  os.Getenv("ENV"),
		HostPort:             os.Getenv("HOST_PORT"),
		FrontendUrl:          os.Getenv("FRONTEND_URL"),
		AuthUrl:              os.Getenv("AUTH_URL"),
		DatabaseUrl:          os.Getenv("DATABASE_URL"),
		SeedDb:               os.Getenv("SEED_DB") == "true",
		CfRegion:             os.Getenv("CF_REGION"),
		CfAccountId:          os.Getenv("CF_ACCOUNT_ID"),
		CfAccessKeyId:        os.Getenv("CF_ACCESS_KEY_ID"),
		CfSecretAccessKey:    os.Getenv("CF_SECRET_ACCESS_KEY"),
		CdnUrl:               os.Getenv("CDN_URL"),
		R2Bucket:             os.Getenv("R2_BUCKET"),
		MeiliSearchHostUrl:   os.Getenv("MEILISEARCH_HOST_URL"),
		MeiliSearchMasterKey: os.Getenv("MEILISEARCH_MASTER_KEY"),
		RedisUrl:             os.Getenv("REDIS_URL"),
		RedisPassword:        os.Getenv("REDIS_PASSWORD"),
	}
}
