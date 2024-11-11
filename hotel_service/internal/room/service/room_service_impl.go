package service

import (
	"errors"
	am "hotel_service/internal/amenity/model"
	ar "hotel_service/internal/amenity/repository"
	mr "hotel_service/internal/room/model"
	rr "hotel_service/internal/room/repository"
	"hotel_service/internal/server/dto"
	"math/big"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RoomServiceImpl struct {
	roomRepository 		rr.RoomRepository
	amenityRepository 	ar.AmenityRepository
}

func NewRoomService(r rr.RoomRepository, a ar.AmenityRepository) RoomService {
	return &RoomServiceImpl{roomRepository: r, amenityRepository: a}
}

func (s *RoomServiceImpl) 	GetByID(id int64) (*dto.RoomResponseDTO, error) {
	room, err := s.roomRepository.GetRoomById(id)
	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Room not found")
		
		return nil, errors.New("Room not found")
	}

	f := new(big.Float).SetRat(&room.Price)

	ams := make([]string, len(room.Amenities))
	for i, amenity := range room.Amenities {
		ams[i] = amenity.Name
	}

	dto := &dto.RoomResponseDTO{
		ID: room.ID,
		Name: room.Name,
		Price: f.String(),
		Amenities: ams,
	}

	return dto, nil
}

func (s *RoomServiceImpl)	CreateRoom(toCreate *dto.RoomRequestDTO, hotel_id int64) (*dto.RoomResponseDTO, error) {
	tx, err := s.roomRepository.Begin()

	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create room")

		return nil, errors.New("Failed to create room")
	}

	defer func() {
		if err != nil {
			s.roomRepository.Rollback(tx)	
		} else {
			logrus.WithTime(time.Now()).Info("Room creation complete!")
			s.roomRepository.Commit(tx)
		}
	}()

	priceBigRat :=	new(big.Rat)
	if _, ok := priceBigRat.SetString(toCreate.Price); !ok {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create room: incorrect price format")

		return nil, errors.New("Failed to create room")
	}

	room := mr.Room{
		Name: toCreate.Name,
		Price: *priceBigRat,
		HotelID: hotel_id,
		Amenities: make([]*am.Amenity, 0),
	}

	for _, amName := range toCreate.Amenities {
		existing, err := s.amenityRepository.GetAmenityIfExists(hotel_id, amName)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				newAmenity := &am.Amenity{
					Name: amName,
					HotelID: hotel_id,
				}
				if err := s.amenityRepository.AddAmenity(tx, newAmenity); err != nil {
					return nil, err
				}
				room.Amenities = append(room.Amenities, newAmenity)

				logrus.WithFields(logrus.Fields{
					"id": newAmenity.ID,
				}).Info("Amenity with name " + amName + "missing, created a new one")
				continue
			}

			return nil, err
		}
		
		room.Amenities = append(room.Amenities, existing)
	}

	logrus.Info("Amenity linkage complete")

	if err := s.roomRepository.AddRoom(tx, &room); err!=nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create room")

		return nil, err
	}

	dto := &dto.RoomResponseDTO{
		ID: room.ID,
		Name: room.Name,
		Price: room.Price.String(),
		Amenities: toCreate.Amenities,
	}

	return dto, nil
}

func (s *RoomServiceImpl)	UpdateRoom(toUpdate *dto.RoomRequestDTO, room_id int64) (*dto.RoomResponseDTO, error) {

}

