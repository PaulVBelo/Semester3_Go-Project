package configs

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

func LoadConfig() {
	if err := godotenv.Load(".env.dev"); err != nil {
		log.Printf("Error loading environment file: %v", err)
	}

	if os.Getenv("KAFKA_ADDRESS") == "" || os.Getenv("KAFKA_TOPIC") == "" ||
		os.Getenv("DELIVERY_SERVICE_ADDRESS") == "" || os.Getenv("KAFKA_GROUP") == "" {
		log.Fatal("Environment variables must be set in .env.dev")
	}
}
