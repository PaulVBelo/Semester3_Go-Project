package service

import "hotel_service/internal/server/dto"

type HotelService interface {
	GetByID(id int64) (dto.HotelDTO, error)
	CreateHotel(hotel *dto.HotelDTO) error
	UpdateHotel(id int64, hotel *dto.HotelDTO) error
}