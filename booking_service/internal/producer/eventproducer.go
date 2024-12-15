package producer

import (
	"booking-service/internal/dto"
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/segmentio/kafka-go"
)

type BookingEventProducer struct {
	writer *kafka.Writer
}

func NewEventProducer(topic string, kafkaAddress string) (*BookingEventProducer, error) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{kafkaAddress},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})

	if w == nil {
		return nil, errors.New("Failed to connect to Kafka Writer")
	}

	return &BookingEventProducer{writer: w}, nil
}

func (ep *BookingEventProducer) Send(event dto.BookingEventDTO) error {
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := kafka.Message{
		Value: data,
	}

	err = ep.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		return err
	}

	log.Printf("Message sent: %s\n", string(data))
	return nil
}

func (ep *BookingEventProducer) Close() error {
	return ep.writer.Close()
}
