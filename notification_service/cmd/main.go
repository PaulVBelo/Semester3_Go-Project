package main

import (
	"log"
	"notification_service/internal/configs"
	"notification_service/internal/kafka"
	"os"
)

func main() {
	configs.LoadConfig()

	address := os.Getenv("KAFKA_ADDRESS")
	topic := os.Getenv("KAFKA_TOPIC")

	kafka.StartKafkaListener(address, topic)

	log.Println("Kafka listener started...")
}
