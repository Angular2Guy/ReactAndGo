package adapter_config

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load("config/properties.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
