package repository

import "hotel_service/internal/hotel/model"

type HotelRepository interface {
	AddHotel(hotel *model.Hotel) error
	GetHotelById(id int64) (*model.Hotel, error)
	UpdateHotel(hotel *model.Hotel) error
	GetAll() ([]model.Hotel, error)
}