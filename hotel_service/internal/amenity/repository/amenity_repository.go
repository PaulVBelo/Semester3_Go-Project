package repository

import (
	"hotel_service/internal/amenity/model"

	"gorm.io/gorm"
)

type AmenityRepository interface { // Этой штуке не нужно уметь начинать/завершать транзакции. Сами по себе услуги не добавляют.
	AddAmenity(tx *gorm.DB, amenity *model.Amenity) error
	AddAmenities(tx *gorm.DB, amenities *[]model.Amenity) error
	// Update нам не понадобится

	GetAmenityById(id int64) (*model.Amenity, error)
	GetAll() ([]model.Amenity, error)
	GetAmenityIfExists(hotel_id int64, name string) (*model.Amenity, error)

	DeleteForRoom(tx *gorm.DB, room_id int64) error
}