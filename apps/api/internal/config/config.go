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
	AwsRegion            string
	AwsAccessKeyId       string
	AwsSecretAccessKey   string
	AwsSessionToken      string
	CdnUrl               string
	S3AssetsBucketName   string
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
		AwsRegion:            os.Getenv("AWS_REGION"),
		AwsAccessKeyId:       os.Getenv("AWS_ACCESS_KEY_ID"),
		AwsSecretAccessKey:   os.Getenv("AWS_SECRET_ACCESS_KEY"),
		AwsSessionToken:      os.Getenv("AWS_SESSION_TOKEN"),
		CdnUrl:               os.Getenv("CDN_URL"),
		S3AssetsBucketName:   os.Getenv("S3_ASSETS_BUCKET_NAME"),
		MeiliSearchHostUrl:   os.Getenv("MEILISEARCH_HOST_URL"),
		MeiliSearchMasterKey: os.Getenv("MEILISEARCH_MASTER_KEY"),
		RedisUrl:             os.Getenv("REDIS_URL"),
		RedisPassword:        os.Getenv("REDIS_PASSWORD"),
	}
}
