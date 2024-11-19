package repository

import (
	"hotel_service/internal/hotel/model"

	"gorm.io/gorm"
)

type HotelRepository interface {
	Begin() (*gorm.DB, error)
	Commit(tx *gorm.DB) error
	Rollback(tx *gorm.DB) error

	AddHotel(tx *gorm.DB, hotel *model.Hotel) error
	UpdateHotel(tx *gorm.DB, hotel *model.Hotel) error

	GetHotelById(id int64) (*model.Hotel, error)
	GetAll() ([]model.Hotel, error)
}