package listener

import (
	"context"
	"encoding/json"
	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"notification_service/internal/handler"
	"notification_service/proto/gen"
	"os"
)

func StartKafkaListener(kafkaAddress, topic string) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaAddress},
		Topic:   topic,
		GroupID: os.Getenv("KAFKA_GROUP"),
	})

	ctx := context.Background()

	logger.WithFields(logrus.Fields{
		"service": "notification_service",
	}).Info("Kafka listener started...")

	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"service": "notification_service",
				"error":   err,
			}).Error("Error reading message from Kafka")
			continue
		}

		bookingEvent := &gen.BookingEvent{}
		if err := json.Unmarshal(m.Value, bookingEvent); err != nil {
			logger.WithFields(logrus.Fields{
				"service": "notification_service",
				"error":   err,
			}).Error("Error unmarshaling Kafka message")
			continue
		}

		logger.WithFields(logrus.Fields{
			"service": "notification_service",
			"event":   bookingEvent,
		}).Info("Decoded Kafka Event")

		if err := handler.HandleBookingEvent(bookingEvent); err != nil {
			logger.WithFields(logrus.Fields{
				"service": "notification_service",
				"error":   err,
			}).Error("Error handling Kafka booking event")
		}
	}
}
