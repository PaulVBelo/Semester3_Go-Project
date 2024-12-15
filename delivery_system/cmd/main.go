package main

import (
	"context"
	"delivery_system/internal/handler"
	"delivery_system/proto/gen"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type server struct {
	gen.UnimplementedDeliverySystemServer
}

func (s *server) SendBooking(_ context.Context, req *gen.BookingEvent) (*gen.BookingResponse, error) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	err := handler.HandleBookingEvent(req)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"service": "delivery_system",
			"error":   err,
		}).Error("failed to handle booking event")
		return nil, err
	}

	return &gen.BookingResponse{
		Success: true,
		Message: "Booking processed successfully",
	}, nil
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.WithFields(logrus.Fields{
			"service": "delivery_system",
			"error":   err,
		}).Fatal("failed to listen")
	}

	s := grpc.NewServer()

	gen.RegisterDeliverySystemServer(s, &server{})

	reflection.Register(s)

	logger.WithFields(logrus.Fields{
		"service": "delivery_system",
	}).Info("Delivery system started listening...")
	if err := s.Serve(lis); err != nil {
		logger.WithFields(logrus.Fields{
			"service": "delivery_system",
			"error":   err,
		}).Fatal("failed to serve")
	}
}
