package config

import (
	"github.com/joho/godotenv"
	"go-sample-rest-api/utils"
)

type Config struct {
	DBUser                 string
	DBPassword             string
	DBName                 string
	DBHost                 string
	DBPort                 string
	JWTSecret              string
	JWTExpirationInSeconds int64
	ServerPort             string
}

var Envs = initConfig()

func initConfig() Config {

	//load env variables
	godotenv.Load()

	return Config{
		DBUser:                 utils.GetEnv("DB_USER", "user"),
		DBPassword:             utils.GetEnv("DB_PASSWORD", "password"),
		DBName:                 utils.GetEnv("DB_NAME", "app"),
		DBHost:                 utils.GetEnv("DB_HOST", "localhost"),
		DBPort:                 utils.GetEnv("DB_PORT", "5432"),
		JWTSecret:              utils.GetEnv("JWT_SECRET", "not-so-secret-now-is-it?"),
		JWTExpirationInSeconds: utils.GetEnvAsInt("JWT_EXPIRATION_IN_SECONDS", 3600*24*7),
		ServerPort:             utils.GetEnv("SERVER_PORT", "8080"),
	}
}
