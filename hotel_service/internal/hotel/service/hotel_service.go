package service

import "hotel_service/internal/server/dto"

type HotelService interface {
	GetByID(id int64) (dto.HotelResponseDTO, error)
	GetAll(id int64) (dto.HotelResponseDTO)
	CreateHotel(hotel *dto.HotelCreateRequestDTO) (dto.HotelResponseDTO, error)
	UpdateHotel(id int64, hotel *dto.HotelCreateRequestDTO) (dto.HotelResponseDTO, error)
}
