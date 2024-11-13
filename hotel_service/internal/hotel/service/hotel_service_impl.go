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
		}).Error("Hotel not found")

		return nil, errors.New("Hotel not found")
	}

	var roomDTOs []*dto.RoomResponseDTO
	for _, room := range hotel.Rooms {
		roomDTO := &dto.RoomResponseDTO{
			ID: room.ID,
			Name: room.Name,
			Price: room.Price.String(),
			Amenities: make([]string, len(room.Amenities)),
		}

		for i, amenity := range room.Amenities {
			roomDTO.Amenities[i] = amenity.Name
		}

		roomDTOs = append(roomDTOs, roomDTO)
	}

	dto := &dto.HotelResponseDTO {
		ID: hotel.ID,
		Name: hotel.Name,
		Adress: hotel.Adress,
		PhoneNumber: hotel.PhoneNumber,
		Rooms: roomDTOs,
	}

	return dto, nil
}

func (s *HotelServiceImpl) GetAll() ([]*dto.HotelResponseDTO, error) {
	hotels, err := s.hotelRepository.GetAll()
	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to retrieve hotels")

		return nil, errors.New("Failed to retrieve hotels")
	}

	hotelDTOs := make([]*dto.HotelResponseDTO, len(hotels))
	for i, hotel := range hotels {
		roomDTOs := make([]*dto.RoomResponseDTO, len(hotel.Rooms))
		for j, room := range hotel.Rooms {
			roomDTOs[j] = &dto.RoomResponseDTO{
				ID: room.ID,
				Name: room.Name,
				Price: room.Price.String(),
				Amenities: make([]string, len(room.Amenities)),
			}
			for k, amenity := range room.Amenities {
				roomDTOs[j].Amenities[k] = amenity.Name
			}
		}

		hotelDTOs[i] = &dto.HotelResponseDTO{
			ID: hotel.ID,
			Name: hotel.Name,
			Adress: hotel.Adress,
			PhoneNumber: hotel.PhoneNumber,
			Rooms: roomDTOs,
		}
	}
	
	return hotelDTOs, nil
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