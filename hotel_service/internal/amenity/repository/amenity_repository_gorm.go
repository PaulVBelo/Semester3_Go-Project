package repository

import (
	"gorm.io/gorm"
	"hotel_service/internal/amenity/model"
)

type amenityRepositoryWithGorm struct {
	db *gorm.DB
}

func NewAmenityRepository(db *gorm.DB) AmenityRepository {
	return &amenityRepositoryWithGorm{db: db}
}

func (r *amenityRepositoryWithGorm) AddAmenity(amenity *model.Amenity) error {
	return r.db.Create(amenity).Error
}

func (r *amenityRepositoryWithGorm) AddAmenities(amenities *[]model.Amenity) error {
	return r.db.Create(amenities).Error
}

func (r *amenityRepositoryWithGorm) GetAmenityById(id int64) (*model.Amenity, error) {
	var amenity model.Amenity
	if err := r.db.First(&amenity, id).Error; err != nil {
		return nil, err
	}
	return &amenity, nil
}

func (r *amenityRepositoryWithGorm) UpdateAmenity(amenity *model.Amenity) error {
	return r.db.Save(amenity).Error
}

func (r *amenityRepositoryWithGorm) GetAll() ([]model.Amenity, error) {
	var amenities []model.Amenity
	result := r.db.Find(&amenities)
	return amenities, result.Error
}