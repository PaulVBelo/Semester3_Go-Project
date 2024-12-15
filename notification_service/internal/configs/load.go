package configs

import (
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
)

func LoadConfig() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	if err := godotenv.Load(".env.dev"); err != nil {
		logger.WithFields(logrus.Fields{
			"service": "notification_service",
			"error":   err,
		}).Error("Error loading environment file")
	}

	if os.Getenv("KAFKA_ADDRESS") == "" || os.Getenv("KAFKA_TOPIC") == "" ||
		os.Getenv("DELIVERY_SERVICE_ADDRESS") == "" || os.Getenv("KAFKA_GROUP") == "" {
		logger.Fatal("Environment variables must be set in .env.dev")
	}
}
