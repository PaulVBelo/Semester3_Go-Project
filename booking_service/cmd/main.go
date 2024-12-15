package main

import (
	"booking-service/internal/config"
	"booking-service/internal/controller"
	"booking-service/internal/db"
	"booking-service/internal/producer"
	"booking-service/internal/router"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"gorm.io/gorm"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig(".env.dev")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Connect to the database
	conn, err := db.ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer closeDB(conn)

	// Initialize BookingEventProducer
	kafkaProducer, err := producer.NewEventProducer("booking-event", "localhost:29093")
	if err != nil {
		log.Fatalf("Error initializing Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// Create controller
	ctrl := controller.NewController(cfg, conn, kafkaProducer)

	// Setup routes
	r := router.NewRouter(ctrl)

	// Setup graceful shutdown
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Port),
		Handler: r,
	}

	go func() {
		log.Printf("Starting server on port %s...", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	gracefulShutdown(server, conn, kafkaProducer)
}

func closeDB(conn *gorm.DB) {
	sqlDB, err := conn.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func gracefulShutdown(server *http.Server, db *gorm.DB, producer *producer.BookingEventProducer) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	closeDB(db)

	// Close Kafka producer
	if err := producer.Close(); err != nil {
		log.Printf("Error closing Kafka producer: %v", err)
	}

	log.Println("Server exited gracefully")
}
