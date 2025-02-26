package utils

import (
	"encoding/json"
	"log"
	"os"

	"github.com/emidiaz3/event-driven-server/models"
	"github.com/joho/godotenv"
)

var ApiKey string
var Config models.Config

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

func InitConfig() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatal("Error leyendo archivo", err)
	}
	err = json.Unmarshal(file, &Config)
	if err != nil {
		log.Fatal("Error convirtiendo archivo", err)
	}
}

func InitApiKey() {
	ApiKey = os.Getenv("API_KEY")
	if ApiKey == "" {
		log.Fatal("API_KEY no está configurada")
	}
}
