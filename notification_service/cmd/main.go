package main

import (
	"github.com/sirupsen/logrus"
	"notification_service/internal/configs"
	"notification_service/internal/listener"
	"os"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	configs.LoadConfig()

	address := os.Getenv("KAFKA_ADDRESS")
	topic := os.Getenv("KAFKA_TOPIC")

	listener.StartKafkaListener(address, topic)

	logger.WithFields(logrus.Fields{
		"service": "notification_svc",
	}).Info("Kafka listener started...")
}
