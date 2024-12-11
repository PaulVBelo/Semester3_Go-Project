package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadConfig() {
	if err := godotenv.Load("../internal/configs/kafka.env"); err != nil {
		log.Printf("Error loading environment file: %v", err)
	}

	if os.Getenv("KAFKA_ADDRESS") == "" || os.Getenv("KAFKA_TOPIC") == "" {
		log.Fatal("KAFKA_ADDRESS and KAFKA_TOPIC must be set in kafka.env")
	}
}
