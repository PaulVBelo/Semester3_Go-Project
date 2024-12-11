package service

import "hotel_service/internal/server/dto"

type HotelService interface {
	GetByID(id int64) (*dto.HotelResponseDTO, error)
	GetAll() ([]*dto.HotelResponseDTO, error)
	CreateHotel(toCreate *dto.HotelCreateRequestDTO) (*dto.HotelShortResponseDTO, error)
	//Create/Update возвращают только факт изменений
	UpdateHotel(id int64, toUpdate *dto.HotelUpdateRequestDTO) (*dto.HotelShortResponseDTO, error)

	// Для Марка <3
	GetExpandedRoomData(id int64) (*dto.FullRoomData, error)
}
