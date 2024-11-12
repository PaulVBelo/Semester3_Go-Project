package service

import (
	"errors"
	ar "hotel_service/internal/amenity/repository"
	hr "hotel_service/internal/hotel/repository"
	rr "hotel_service/internal/room/repository"
	"hotel_service/internal/server/dto"
	"time"

	"github.com/sirupsen/logrus"
)

type HotelServiceImpl struct {
	hotelRepository   hr.HotelRepository
	roomRepository    rr.RoomRepository
	amenityRepository ar.AmenityRepository
}

func NewHotelService(r rr.RoomRepository, a ar.AmenityRepository, h hr.HotelRepository) HotelService {
	return &HotelServiceImpl{hotelRepository: h, roomRepository: r, amenityRepository: a}
}

func (s *HotelServiceImpl) GetByID(id int64) (*dto.HotelResponseDTO, error) {
	hotel, err := s.hotelRepository.GetHotelById(id)
	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Room not found")

		return nil, errors.New("Room not found")
	}


}

func (s *HotelServiceImpl) GetAll(id int64) (*dto.HotelResponseDTO, error) {

}

func (s *HotelServiceImpl) CreateHotel(hotel *dto.HotelCreateRequestDTO) (*dto.HotelResponseDTO, error) {

}

func (s *HotelServiceImpl) UpdateHotel(id int64, hotel *dto.HotelCreateRequestDTO) (*dto.HotelResponseDTO, error) {
	// Здесь упрощу логику - не буду обновлять трёхмерную развёртку, только сам отель
}

func checkUnique(strings []string) ([]string, bool) { // DupeChecker
	stringMap := make(map[string]int)
	var duplicates []string
	for _, str := range strings {
		stringMap[str]++
		if stringMap[str] == 2 {
			duplicates = append(duplicates, str)
		}
	}
	isUnique := len(duplicates) == 0
	return duplicates, isUnique
}