package service

import (
	"errors"
	am "hotel_service/internal/amenity/model"
	ar "hotel_service/internal/amenity/repository"
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

type RoomServiceImpl struct {
	roomRepository    rr.RoomRepository
	amenityRepository ar.AmenityRepository
}

func NewRoomService(r rr.RoomRepository, a ar.AmenityRepository) RoomService {
	return &RoomServiceImpl{roomRepository: r, amenityRepository: a}
}

func (s *RoomServiceImpl) GetByID(id int64) (*dto.RoomResponseDTO, error) {
	room, err := s.roomRepository.GetRoomById(id)
	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to retrieve room from database")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &se.NotFoundError{"Room not found"}
		}
		return nil, &se.InternalServerError{"Failed to retrieve room"}
		
	}

	f := new(big.Float).SetRat(room.Price.Get())

	ams := make([]string, len(room.Amenities))
	for i, amenity := range room.Amenities {
		ams[i] = amenity.Name
	}

	dto := &dto.RoomResponseDTO{
		ID:        room.ID,
		Name:      room.Name,
		Price:     f.String(),
		Amenities: ams,
	}

	return dto, nil
}

func (s *RoomServiceImpl) CreateRoom(toCreate *dto.RoomCreateRequestDTO, hotel_id int64) (*dto.RoomResponseDTO, error) {
	tx, err := s.roomRepository.Begin()

	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create room")

		return nil, &se.InternalServerError{"Failed to create room"}
	}

	defer func() {
		if err != nil {
			s.roomRepository.Rollback(tx)
		} else {
			logrus.WithTime(time.Now()).Info("Room creation complete!")
			s.roomRepository.Commit(tx)
		}
	}()
	
	priceBigRat := rm.BigRat{Rat: new(big.Rat)}
	if _, ok := priceBigRat.SetString(toCreate.Price); !ok {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create room: incorrect price format")

		return nil, &se.BadRequestError{"Incorrect price format"}
	}

	room := rm.Room{
		Name:      toCreate.Name,
		Price:     priceBigRat,
		HotelID:   hotel_id,
		Amenities: make([]*am.Amenity, 0),
	}

	dupes, ok := checkUnique(toCreate.Amenities)

	if !ok {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"dupes": dupes,
		}).Error("Duplicate amenity")

		return nil, &se.BadRequestError{"Duplicate amenities: " + strings.Join(dupes, ", ")}
	}

	for _, amName := range toCreate.Amenities {
		existing, err := s.amenityRepository.GetAmenityIfExists(hotel_id, amName)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newAmenity := &am.Amenity{
					Name:    amName,
					HotelID: hotel_id,
				}
				if err := s.amenityRepository.AddAmenity(tx, newAmenity); err != nil {
					return nil,  &se.InternalServerError{"Failed to create room"}
				}
				room.Amenities = append(room.Amenities, newAmenity)

				logrus.WithFields(logrus.Fields{
					"id": newAmenity.ID,
					"name": amName,
				}).Debug("Amenity with name missing, created a new one")
				continue
			}

			return nil, &se.InternalServerError{"Failed to create room"}
		}

		room.Amenities = append(room.Amenities, existing)
	}

	logrus.Trace("Amenity linkage complete")

	if err := s.roomRepository.AddRoom(tx, &room); err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create room")

		if (errors.As(err, &se.BadRequestError{})) {
			return nil, err
		}

		if (errors.Is(err, gorm.ErrCheckConstraintViolated)) {
			return nil, &se.BadRequestError{Message: "Bad request"}
		}
		
		return nil, &se.InternalServerError{"Failed to create room"}
	}

	dto := &dto.RoomResponseDTO{
		ID:        room.ID,
		Name:      room.Name,
		Price:     room.Price.String(),
		Amenities: toCreate.Amenities,
	}

	return dto, nil
}

func (s *RoomServiceImpl) UpdateRoom(toUpdate *dto.RoomUpdateRequestDTO, room_id int64) (*dto.RoomResponseDTO, error) {
	tx, err := s.roomRepository.Begin()

	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to update room")

		return nil, &se.InternalServerError{Message: "Failed to update room"}
	}

	defer func() {
		if err != nil {
			s.roomRepository.Rollback(tx)
		} else {
			logrus.WithTime(time.Now()).Info("Room update complete!")
			s.roomRepository.Commit(tx)
		}
	}()

	room, err := s.roomRepository.GetRoomById(room_id)
	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to update room: room not found")
		return nil, &se.NotFoundError{"Room not found"}
	}

	if (toUpdate.Name != nil) {
		room.Name = *toUpdate.Name
	}

	if (toUpdate.Price != nil) {
		if _, ok := room.Price.SetString(*toUpdate.Price); !ok {
			logrus.WithTime(time.Now()).WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Failed to update room: incorrect price format")

			return nil, &se.BadRequestError{"Incorrect price format"}
		}
	}

	if len(toUpdate.Amenities) > 0 {
		dupes, ok := checkUnique(toUpdate.Amenities)

		if !ok {
			logrus.WithTime(time.Now()).WithFields(logrus.Fields{
				"dupes": dupes,
			}).Error("Duplicate amenity")

			return nil, &se.BadRequestError{"Duplicate amenities: " + strings.Join(dupes, ", ")}
		}
		room.Amenities = make([]*am.Amenity, 0)

		for _, amName := range toUpdate.Amenities {
			existing, err := s.amenityRepository.GetAmenityIfExists(room.HotelID, amName)
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					newAmenity := &am.Amenity{
						Name:    amName,
						HotelID: room.HotelID,
					}
					if err := s.amenityRepository.AddAmenity(tx, newAmenity); err != nil {
						return nil, &se.InternalServerError{Message: "Failed to update room"}
					}
					room.Amenities = append(room.Amenities, newAmenity)
	
					logrus.WithFields(logrus.Fields{
						"id": newAmenity.ID,
						"name": amName,
					}).Debug("Amenity with name missing, created a new one")
					continue
				}
	
				return nil, &se.InternalServerError{Message: "Failed to update room"}
			}
	
			room.Amenities = append(room.Amenities, existing)
		}
	}
	logrus.Trace("Amenity linkage complete")

	if err := s.roomRepository.UpdateRoom(tx, room); err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to update room")

		if (errors.As(err, &se.BadRequestError{})) {
			return nil, err
		}

		if (errors.Is(err, gorm.ErrCheckConstraintViolated)) {
			return nil, &se.BadRequestError{Message: "Bad request"}
		}

		return nil, &se.InternalServerError{Message: "Failed to update room"}
	}

	dto := &dto.RoomResponseDTO{
		ID: room.ID,
		Name: room.Name,
		Price: room.Price.String(),
		Amenities: toUpdate.Amenities,
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