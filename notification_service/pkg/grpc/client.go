package grpc

import (
	"context"
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"
	pb "notification_service/proto/gen"
)

type DeliverySystemClient struct {
	client pb.DeliverySystemClient
	conn   *grpc.ClientConn
	mu     sync.Mutex
}

func NewDeliverySystemClient(address string) (*DeliverySystemClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("Failed to connect to DeliverySystem gRPC server: %v", err)
		return nil, err
	}

	return &DeliverySystemClient{
		client: pb.NewDeliverySystemClient(conn),
		conn:   conn,
	}, nil
}

func (d *DeliverySystemClient) SendBooking(ctx context.Context, req *pb.BookingEvent) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	d.mu.Lock()
	defer d.mu.Unlock()

	res, err := d.client.SendBooking(ctx, req)
	if err != nil {
		log.Printf("Error sending booking to DeliverySystem: %v", err)
		return err
	}

	if res.Success {
		log.Printf("Successfully sent booking to DeliverySystem: %v", res.Message)
	} else {
		log.Printf("DeliverySystem responded with failure: %v", res.Message)
	}

	return nil
}

func (d *DeliverySystemClient) Close() {
	if err := d.conn.Close(); err != nil {
		log.Printf("Failed to close gRPC connection: %v", err)
	}
}
