package main

import (
	"booking-service/internal/api"
	"booking-service/internal/config"
	"booking-service/internal/db"
	"fmt"
	"log"
	"net/http"
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
	defer conn.Close()

	// Создание маршрутов
	router := api.NewRouter(conn)

	// Запуск сервиса
	log.Printf("Starting server on port %s...", cfg.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), router))
}
