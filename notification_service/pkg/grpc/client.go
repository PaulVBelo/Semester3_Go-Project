package grpc

import (
	"context"
	"log"
	"os"
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
	stderrLogger := log.New(os.Stderr, "", log.Ldate|log.Ltime)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		stderrLogger.Printf("Failed to connect to DeliverySystem handler: %v", err)
		return nil, err
	}

	return &DeliverySystemClient{
		client: gen.NewDeliverySystemClient(conn),
		conn:   conn,
	}, nil
}

func (d *DeliverySystemClient) SendBooking(ctx context.Context, event *gen.BookingEvent) error {
	stderrLogger := log.New(os.Stderr, "", log.Ldate|log.Ltime)

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	d.mu.Lock()
	defer d.mu.Unlock()

	res, err := d.client.SendBooking(ctx, event)
	if err != nil {
		stderrLogger.Printf("Error sending booking to DeliverySystem: %v", err)
		return err
	}

	if res.Success {
		log.Printf("Successfully sent booking to DeliverySystem: %v", res.Message)
	} else {
		stderrLogger.Printf("DeliverySystem responded with failure: %v", res.Message)
	}

	return nil
}

func (d *DeliverySystemClient) Close() {
	stderrLogger := log.New(os.Stderr, "", log.Ldate|log.Ltime)

	if err := d.conn.Close(); err != nil {
		stderrLogger.Printf("Failed to close gRPC connection: %v", err)
	}
}
