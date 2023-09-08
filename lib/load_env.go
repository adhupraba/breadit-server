package lib

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type envConfig struct {
	Port               string
	DbUrl              string
	JwtSecret          string
	RedisUrl           string
	Env                string
	CorsAllowedOrigins []string
}

var EnvConfig envConfig

func LoadEnv() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Unable to load .env:", err)
	}

	var env string = os.Getenv("ENV")

	if env == "" {
		env = "development"
	}

	EnvConfig = envConfig{
		Port:               os.Getenv("PORT"),
		DbUrl:              os.Getenv("DB_URL"),
		JwtSecret:          os.Getenv("JWT_SECRET"),
		RedisUrl:           os.Getenv("REDIS_URL"),
		Env:                env,
		CorsAllowedOrigins: []string{},
	}

	if EnvConfig.Env == "" {
		EnvConfig.Env = "development"
	}

	if os.Getenv("CORS_ALLOWED_ORIGINGS") == "" {
		EnvConfig.CorsAllowedOrigins = []string{"*"}
	} else {
		EnvConfig.CorsAllowedOrigins = strings.Split(os.Getenv("CORS_ALLOWED_ORIGINGS"), " ")
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
