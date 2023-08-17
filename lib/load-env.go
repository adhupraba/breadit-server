package lib

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type envConfig struct {
	Port      string
	DbUrl     string
	JwtSecret string
	RedisUrl  string
}

var EnvConfig envConfig

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Unable to load .env:", err)
	}

	EnvConfig = envConfig{
		Port:      os.Getenv("PORT"),
		DbUrl:     os.Getenv("DB_URL"),
		JwtSecret: os.Getenv("JWT_SECRET"),
		RedisUrl:  os.Getenv("REDIS_URL"),
	}

	if EnvConfig.Port == "" {
		log.Fatal("PORT is not found in the environment")
	}

	if EnvConfig.DbUrl == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	if EnvConfig.JwtSecret == "" {
		log.Fatal("JWT_SECRET is not found in the environment")
	}

	if EnvConfig.RedisUrl == "" {
		log.Fatal("REDIS_URL is not found in the environment")
	}
}
