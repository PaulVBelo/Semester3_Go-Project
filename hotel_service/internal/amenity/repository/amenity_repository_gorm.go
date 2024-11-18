package repository

import (
	"hotel_service/internal/amenity/model"

	"gorm.io/gorm"
)

type AmenityRepositoryWithGorm struct {
	db *gorm.DB
}

func NewAmenityRepository(db *gorm.DB) AmenityRepository {
	return &AmenityRepositoryWithGorm{db: db}
}

func (r *AmenityRepositoryWithGorm) AddAmenity(tx *gorm.DB, amenity *model.Amenity) error {
	return tx.Create(amenity).Error
}

func (r *AmenityRepositoryWithGorm) AddAmenities(tx *gorm.DB, amenities *[]model.Amenity) error {
	return tx.Create(amenities).Error
}

func (r *AmenityRepositoryWithGorm) GetAmenityById(id int64) (*model.Amenity, error) {
	var amenity model.Amenity
	if err := r.db.First(&amenity, id).Error; err != nil {
		return nil, err
	}
	return &amenity, nil
}

func (r *AmenityRepositoryWithGorm) GetAll() ([]model.Amenity, error) {
	var amenities []model.Amenity
	result := r.db.Find(&amenities)
	return amenities, result.Error
}

func (r *AmenityRepositoryWithGorm) GetAmenityIfExists(hotel_id int64, name string) (*model.Amenity, error) {
	var amenity model.Amenity
	if err := r.db.Where("amenity_name = ? AND hotel_id = ?", name, hotel_id).First(&amenity).Error; err != nil {
		return nil, err
	}
	return &amenity, nil
}

func (r *AmenityRepositoryWithGorm) DeleteForRoom(tx *gorm.DB, room_id int64) error {
	return tx.Exec("DELETE FROM room_x_amenity WHERE room_id = ?", room_id).Error
}
