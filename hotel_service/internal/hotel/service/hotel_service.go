package service

import "hotel_service/internal/server/dto"

type HotelService interface {
	GetByID(id int64) (dto.HotelRequestDTO, error)
	CreateHotel(hotel *dto.HotelRequestDTO) error
	UpdateHotel(id int64, hotel *dto.HotelRequestDTO) error
}
