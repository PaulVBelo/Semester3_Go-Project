package grpc

import (
	"context"
	"github.com/sirupsen/logrus"
	"sync"
	"time"

	"google.golang.org/grpc"
	"notification_service/proto/gen"
)

type DeliverySystemClient struct {
	client gen.DeliverySystemClient
	conn   *grpc.ClientConn
	mu     sync.Mutex
}

func NewDeliverySystemClient(address string) (*DeliverySystemClient, error) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		logger.WithFields(logrus.Fields{
			"service": "notification_service",
			"address": address,
			"error":   err,
		}).Error("Failed to connect to DeliverySystem handler")
		return nil, err
	}

	return &DeliverySystemClient{
		client: gen.NewDeliverySystemClient(conn),
		conn:   conn,
	}, nil
}

func (d *DeliverySystemClient) SendBooking(ctx context.Context, event *gen.BookingEvent) error {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	d.mu.Lock()
	defer d.mu.Unlock()

	res, err := d.client.SendBooking(ctx, event)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"service": "notification_service",
			"error":   err,
		}).Error("Error sending booking to DeliverySystem")
		return err
	}

	if res.Success {
		logger.WithFields(logrus.Fields{
			"service": "notification_service",
			"message": res.Message,
		}).Info("Successfully sent booking to DeliverySystem")
	} else {
		logger.WithFields(logrus.Fields{
			"service": "notification_service",
			"message": res.Message,
		}).Error("DeliverySystem responded with failure")
	}

	return nil
}

func (d *DeliverySystemClient) Close() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	if err := d.conn.Close(); err != nil {
		logger.WithFields(logrus.Fields{
			"service": "notification_service",
			"error":   err,
		}).Error("Failed to close gRPC connection")
	}
}
