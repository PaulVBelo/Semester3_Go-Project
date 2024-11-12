package repository

import (
	"hotel_service/internal/hotel/model"

	"gorm.io/gorm"
)

type HotelRepositoryWithGorm struct {
	db *gorm.DB
}

func NewHotelRepository(db *gorm.DB) HotelRepository {
	return &HotelRepositoryWithGorm{db: db}
}

func (r *HotelRepositoryWithGorm) Begin() (*gorm.DB, error) {
	return r.db.Begin(), nil
}

func (r *HotelRepositoryWithGorm) Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (r *HotelRepositoryWithGorm) Rollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func (r *HotelRepositoryWithGorm) AddHotel(hotel *model.Hotel) error {
	return r.db.Create(hotel).Error
}

func (r *HotelRepositoryWithGorm) GetHotelById(id int64) (*model.Hotel, error) {
	var hotel model.Hotel
	if err := r.db.Preload("Rooms").First(&hotel, id).Error; err != nil {
		return nil, err
	}
	return &hotel, nil
}

func (r *HotelRepositoryWithGorm) UpdateHotel(hotel *model.Hotel) error {
	return r.db.Save(hotel).Error
}

func (r *HotelRepositoryWithGorm) GetAll() ([]model.Hotel, error) {
	var hotels []model.Hotel
	result := r.db.Preload("Rooms").Find(&hotels)
	return hotels, result.Error
}
