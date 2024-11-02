package repository

import (
	"gorm.io/gorm"
	"hotel_service/internal/hotel/model"
)

type hotelRepositoryWithGorm struct {
	db *gorm.DB
}

func NewHotelRepository(db *gorm.DB) HotelRepository {
	return &hotelRepositoryWithGorm{db: db}
}

func (r *hotelRepositoryWithGorm) AddHotel(hotel *model.Hotel) error {
	return r.db.Create(hotel).Error
}

func (r *hotelRepositoryWithGorm) GetHotelById(id int64) (*model.Hotel, error) {
	var hotel model.Hotel
	if err := r.db.Preload("Rooms").First(&hotel, id).Error; err != nil {
		return nil, err
	}
	return &hotel, nil
}

func (r *hotelRepositoryWithGorm) UpdateHotel(hotel *model.Hotel) error {
	return r.db.Save(hotel).Error
}

func (r *hotelRepositoryWithGorm) GetAll() ([]model.Hotel, error) {
	var hotels []model.Hotel
	result := r.db.Preload("Rooms").Find(&hotels)
	return hotels, result.Error
}