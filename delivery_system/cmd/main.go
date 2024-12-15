package main

import (
	"context"
	"delivery_system/internal/handler"
	"delivery_system/proto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type server struct {
	gen.UnimplementedDeliverySystemServer
}

func (s *server) SendBooking(_ context.Context, req *gen.BookingEvent) (*gen.BookingResponse, error) {

	handler.HandleBookingEvent(req)

	return &gen.BookingResponse{
		Success: true,
		Message: "Booking processed successfully",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	gen.RegisterDeliverySystemServer(s, &server{})

	reflection.Register(s)

	log.Println("Delivery system started listening...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}