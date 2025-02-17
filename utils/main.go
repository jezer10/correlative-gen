package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var ApiKey string

func InitEnv() {
	env := os.Getenv("APP_ENV")
	var envFile string
	if env == "production" {
		envFile = ".env.production"
	} else {
		envFile = ".env.development"
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Fatal("Error loading .env file")
	}

}

func InitApiKey() {
	ApiKey = os.Getenv("API_KEY")
	if ApiKey == "" {
		log.Fatal("API_KEY no está configurada")
	}
}
