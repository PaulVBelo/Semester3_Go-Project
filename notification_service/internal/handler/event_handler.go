package handler

import (
	"context"
	"github.com/sirupsen/logrus"
	"notification_service/pkg/grpc"
	"notification_service/proto/gen"
	"os"
)

func HandleBookingEvent(event *gen.BookingEvent) error {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.WithFields(logrus.Fields{
		"service": "notification_service",
		"event":   event,
	}).Info("Received booking event")

	client, err := grpc.NewDeliverySystemClient(os.Getenv("DELIVERY_SERVICE_ADDRESS"))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"service": "notification_service",
			"error":   err,
		}).Error("Failed to create delivery client")
		return err
	}
	defer client.Close()

	ctx := context.Background()
	err = client.SendBooking(ctx, event)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"service": "notification_service",
			"error":   err,
		}).Error("Failed to send booking event to delivery system")
		return err
	}

	return nil
}
