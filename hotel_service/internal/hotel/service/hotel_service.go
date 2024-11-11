package service

import "hotel_service/internal/server/dto"

type HotelService interface {
	GetByID(id int64) (dto.HotelCreateRequestDTO, error)
	CreateHotel(hotel *dto.HotelCreateRequestDTO) error
	UpdateHotel(id int64, hotel *dto.HotelCreateRequestDTO) error
}
