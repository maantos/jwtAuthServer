package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVarables() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env from file")
	}
}
