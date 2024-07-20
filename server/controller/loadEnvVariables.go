package controller

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	// provide .env file's location in following function
	er := godotenv.Load("../.env")

	if er != nil {
		log.Fatalf("Error loading .env file %v", er)
	}
}
