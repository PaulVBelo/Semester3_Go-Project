package service

import (
	ar "hotel_service/internal/amenity/repository"
	hr "hotel_service/internal/hotel/repository"
	rr "hotel_service/internal/room/repository"
	"hotel_service/internal/server/dto"
)

type HotelServiceImpl struct {
	hotelRepository   hr.HotelRepository
	roomRepository    rr.RoomRepository
	amenityRepository ar.AmenityRepository
}

func NewHotelService(r rr.RoomRepository, a ar.AmenityRepository, h hr.HotelRepository) HotelService {
	return &HotelServiceImpl{hotelRepository: h, roomRepository: r, amenityRepository: a}
}

func (s *HotelServiceImpl) GetByID(id int64) (dto.HotelResponseDTO, error) {

}

func (s *HotelServiceImpl) GetAll(id int64) (dto.HotelResponseDTO) {

}

func (s *HotelServiceImpl) CreateHotel(hotel *dto.HotelCreateRequestDTO) (dto.HotelResponseDTO, error) {

}

func (s *HotelServiceImpl) UpdateHotel(id int64, hotel *dto.HotelCreateRequestDTO) (dto.HotelResponseDTO, error) {

}