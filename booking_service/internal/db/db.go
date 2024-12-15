package db

import (
	"booking-service/internal/config"
	"booking-service/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}

	log.Println("Connected to the database successfully!")

	// Автоматическая миграция
	err = db.AutoMigrate(&models.Booking{})
	if err != nil {
		return nil, fmt.Errorf("error during migration: %v", err)
	}

	return db, nil
}
