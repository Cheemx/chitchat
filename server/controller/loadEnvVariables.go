package controller

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	er := godotenv.Load()

	if er != nil {
		log.Fatalf("Error loading .env file %v", er)
	}
}
