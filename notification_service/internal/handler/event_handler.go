package handler

import (
	"context"
	"log"
	"notification_service/pkg/grpc"
	"notification_service/proto/gen"
	"os"
)

func HandleBookingEvent(event *gen.BookingEvent) error {
	log.Printf("Received booking event: %+v\n", event)

	client, err := grpc.NewDeliverySystemClient(os.Getenv("DELIVERY_SERVICE_ADDRESS"))
	if err != nil {
		log.Printf("Failed to create delivery client: %v", err)
		return err
	}
	defer client.Close()

	ctx := context.Background()
	err = client.SendBooking(ctx, event)
	if err != nil {
		log.Printf("Failed to send booking event to delivery system: %v", err)
		return err
	}

	return nil
}
