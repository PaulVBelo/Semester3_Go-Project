package repository

import "hotel_service/internal/amenity/model"

type AmenityRepository interface {
	AddAmenity(amenity *model.Amenity) error
	AddAmenities(amenities *[]model.Amenity) error
	GetAmenityById(id int64) (*model.Amenity, error)
	UpdateAmenity(amenity *model.Amenity) error
	GetAll() ([]model.Amenity, error)
}