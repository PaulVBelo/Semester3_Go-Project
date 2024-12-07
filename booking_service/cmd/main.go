package main

import (
	"booking-service/internal/config"
	"booking-service/internal/controller"
	"booking-service/internal/db"
	"booking-service/internal/router"
	"context"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig(".env.dev")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Подключение к базе данных
	conn, err := db.ConnectPostgres(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer closeDB(conn)

	// Создаем контроллер
	ctrl := controller.NewController(cfg, conn)

	// Настройка маршрутов
	r := router.NewRouter(ctrl)

	// Настройка graceful shutdown
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

	gracefulShutdown(server, conn)
}

func closeDB(conn *gorm.DB) {
	sqlDB, err := conn.DB()
	if err == nil {
		sqlDB.Close()
	}
}

func gracefulShutdown(server *http.Server, db *gorm.DB) {
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
	log.Println("Server exited gracefully")
}
