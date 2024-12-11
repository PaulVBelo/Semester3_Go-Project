package kafka

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"log"
	"notification_service/internal/handler"
	"notification_service/proto/gen"
)

func StartKafkaListener(kafkaAddress, topic string) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaAddress},
		Topic:   topic,
		GroupID: "notification-service-group",
	})

	ctx := context.Background()

	log.Println("Kafka listener started...")
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Printf("error reading message from Kafka: %v\n", err)
			continue
		}

		bookingEvent := &gen.BookingEvent{}
		if err := json.Unmarshal(m.Value, bookingEvent); err != nil {
			log.Printf("error unmarshaling Kafka message: %v\n", err)
			continue
		}

		log.Printf("Decoded Kafka Event: %+v\n", bookingEvent)

		if err := handler.HandleBookingEvent(bookingEvent); err != nil {
			log.Printf("error handling Kafka booking event: %v", err)
		}
	}
}
