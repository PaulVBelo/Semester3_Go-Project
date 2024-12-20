package service

import (
	"errors"
	am "hotel_service/internal/amenity/model"
	ar "hotel_service/internal/amenity/repository"
	hm "hotel_service/internal/hotel/model"
	hr "hotel_service/internal/hotel/repository"
	rm "hotel_service/internal/room/model"
	rr "hotel_service/internal/room/repository"
	"hotel_service/internal/server/dto"
	se "hotel_service/internal/server/errors"
	"math/big"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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
		}).Error("Unable to retrieve hotel")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &se.NotFoundError{"Hotel not found"}
		}

		return nil, errors.New("Failed to retrieve hotel")
	}

	var roomDTOs []*dto.RoomResponseDTO
	for _, room := range hotel.Rooms {
		roomDTO := &dto.RoomResponseDTO{
			ID:        room.ID,
			Name:      room.Name,
			Price:     room.Price.String(),
			Amenities: make([]string, len(room.Amenities)),
		}

		for i, amenity := range room.Amenities {
			roomDTO.Amenities[i] = amenity.Name
		}

		roomDTOs = append(roomDTOs, roomDTO)
	}

	dto := &dto.HotelResponseDTO{
		ID:          hotel.ID,
		Name:        hotel.Name,
		Adress:      hotel.Adress,
		PhoneNumber: hotel.PhoneNumber,
		Rooms:       roomDTOs,
	}

	return dto, nil
}

func (s *HotelServiceImpl) GetExpandedRoomData(id int64) (*dto.FullRoomData, error) {
	room, err := s.roomRepository.GetRoomById(id)
	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Unable to retrieve room")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &se.NotFoundError{"Room not found"}
		}

		return nil, errors.New("Failed to retrieve room")
	}

	hotel, err := s.hotelRepository.GetHotelById(room.HotelID)
	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Unable to retrieve hotel")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &se.NotFoundError{"Hotel not found"}
		}

		return nil, errors.New("Failed to retrieve hotel")
	}

	responseDTO := dto.FullRoomData{
		ID:          hotel.ID,
		Name:        hotel.Name,
		Adress:      hotel.Adress,
		PhoneNumber: hotel.PhoneNumber,
		RoomId:      room.ID,
		RoomName:    room.Name,
		Price:       room.Price.String(),
		Amenities:   make([]string, len(room.Amenities)),
	}

	for i, amenity := range room.Amenities {
		responseDTO.Amenities[i] = amenity.Name
	}

	return &responseDTO, nil
}

func (s *HotelServiceImpl) GetAll() ([]*dto.HotelResponseDTO, error) {
	hotels, err := s.hotelRepository.GetAll()
	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to retrieve hotels")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &se.NotFoundError{"No hotels found"}
		}

		return nil, errors.New("Failed to retrieve hotels")
	}

	hotelDTOs := make([]*dto.HotelResponseDTO, len(hotels))
	for i, hotel := range hotels {
		roomDTOs := make([]*dto.RoomResponseDTO, len(hotel.Rooms))
		for j, room := range hotel.Rooms {
			roomDTOs[j] = &dto.RoomResponseDTO{
				ID:        room.ID,
				Name:      room.Name,
				Price:     room.Price.String(),
				Amenities: make([]string, len(room.Amenities)),
			}
			for k, amenity := range room.Amenities {
				roomDTOs[j].Amenities[k] = amenity.Name
			}
		}

		hotelDTOs[i] = &dto.HotelResponseDTO{
			ID:          hotel.ID,
			Name:        hotel.Name,
			Adress:      hotel.Adress,
			PhoneNumber: hotel.PhoneNumber,
			Rooms:       roomDTOs,
		}
	}

	return hotelDTOs, nil
}

func (s *HotelServiceImpl) CreateHotel(toCreate *dto.HotelCreateRequestDTO) (*dto.HotelShortResponseDTO, error) {
	// Это ассимптотический цирк с конями
	tx, err := s.hotelRepository.Begin()
	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create hotel")

		return nil, &se.InternalServerError{"Failed to create hotel"}
	}

	defer func() {
		if err != nil {
			s.hotelRepository.Rollback(tx)
		} else {
			logrus.WithTime(time.Now()).Info("Hotel creation complete")
			s.hotelRepository.Commit(tx)
		}
	}()

	hotel := &hm.Hotel{
		Name:        toCreate.Name,
		Adress:      toCreate.Adress,
		PhoneNumber: toCreate.PhoneNumber,
		Rooms:       make([]*rm.Room, 0),
	}

	if err := s.hotelRepository.AddHotel(tx, hotel); err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create hotel")

		if (errors.As(err, &se.BadRequestError{})) {
			return nil, err
		}

		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			return nil, &se.BadRequestError{Message: "Bad request"}
		}

		return nil, &se.InternalServerError{"Failed to create hotel"}
	}

	if len(toCreate.Rooms) == 0 {
		return &dto.HotelShortResponseDTO{
			ID:          hotel.ID,
			Name:        hotel.Name,
			Adress:      hotel.Adress,
			PhoneNumber: hotel.PhoneNumber,
		}, nil
	}

	amenityCache := make(map[string]*am.Amenity)

	for _, roomCreateDTO := range toCreate.Rooms {

		priceBigRat := rm.BigRat{new(big.Rat)}
		if _, ok := priceBigRat.SetString(roomCreateDTO.Price); !ok {
			logrus.WithTime(time.Now()).WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Failed to create room: incorrect price format")

			return nil, &se.BadRequestError{"Failed to create room - incorrect price format: " + roomCreateDTO.Price}
		}

		room := &rm.Room{
			Name:      roomCreateDTO.Name,
			Price:     priceBigRat,
			HotelID:   hotel.ID,
			Amenities: make([]*am.Amenity, 0),
		}

		dupes, ok := checkUnique(roomCreateDTO.Amenities)

		if !ok {
			logrus.WithTime(time.Now()).WithFields(logrus.Fields{
				"dupes": dupes,
			}).Error("Duplicate amenities")

			return nil, &se.BadRequestError{"Duplicate amenities: " + strings.Join(dupes, ", ") + " at " + roomCreateDTO.Name}
		}

		for _, amName := range roomCreateDTO.Amenities {

			if existing, found := amenityCache[amName]; found {
				room.Amenities = append(room.Amenities, existing)
				continue
			}

			newAmenity := &am.Amenity{
				Name:    amName,
				HotelID: hotel.ID,
			}

			if err := s.amenityRepository.AddAmenity(tx, newAmenity); err != nil {
				logrus.WithTime(time.Now()).WithFields(logrus.Fields{
					"error": err.Error(),
				}).Error("Failed to create hotel: to create amenity " + amName)

				return nil, &se.InternalServerError{"Failed to create hotel: to create amenity " + amName}
			}

			room.Amenities = append(room.Amenities, newAmenity)
			amenityCache[amName] = newAmenity
		}

		if err := s.roomRepository.AddRoom(tx, room); err != nil {
			logrus.WithTime(time.Now()).WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Failed to create hotel: creating room " + roomCreateDTO.Name)

			if (errors.As(err, &se.BadRequestError{})) {
				return nil, err
			}

			if errors.Is(err, gorm.ErrCheckConstraintViolated) {
				return nil, &se.BadRequestError{Message: "Bad request"}
			}

			return nil, &se.InternalServerError{"Failed to create hotel: creating room " + roomCreateDTO.Name}
		}

		hotel.Rooms = append(hotel.Rooms, room)
	}

	if err := s.hotelRepository.UpdateHotel(tx, hotel); err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create hotel with rooms")

		return nil, &se.InternalServerError{"Failed to create hotel with rooms"}
	}

	dto := &dto.HotelShortResponseDTO{
		ID:          hotel.ID,
		Name:        hotel.Name,
		Adress:      hotel.Adress,
		PhoneNumber: hotel.PhoneNumber,
	}

	return dto, nil
}

func (s *HotelServiceImpl) UpdateHotel(id int64, toUpdate *dto.HotelUpdateRequestDTO) (*dto.HotelShortResponseDTO, error) {
	tx, err := s.hotelRepository.Begin()

	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to update hotel")

		return nil, &se.InternalServerError{"Failed to update hotel"}
	}

	defer func() {
		if err != nil {
			s.roomRepository.Rollback(tx)
		} else {
			logrus.WithTime(time.Now()).Info("Hotel update complete!")
			s.roomRepository.Commit(tx)
		}
	}()

	hotel, err := s.hotelRepository.GetHotelById(id)
	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Unable to retrieve hotel")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &se.NotFoundError{"Hotel not found"}
		}

		return nil, &se.InternalServerError{"Hotel not found"}
	}

	if toUpdate.Name != nil {
		hotel.Name = *toUpdate.Name
	}

	if toUpdate.Adress != nil {
		hotel.Adress = *toUpdate.Adress
	}

	if toUpdate.PhoneNumber != nil {
		hotel.PhoneNumber = *toUpdate.PhoneNumber
	}

	if err := s.hotelRepository.UpdateHotel(tx, hotel); err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to update hotel")

		if (errors.As(err, &se.BadRequestError{})) {
			return nil, err
		}

		if errors.Is(err, gorm.ErrCheckConstraintViolated) {
			return nil, &se.BadRequestError{Message: "Bad request"}
		}

		return nil, &se.InternalServerError{"Failed to update hotel"}
	}

	dto := &dto.HotelShortResponseDTO{
		ID:          hotel.ID,
		Name:        hotel.Name,
		Adress:      hotel.Adress,
		PhoneNumber: hotel.PhoneNumber,
	}

	return dto, nil
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
