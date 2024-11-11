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
			s.roomRepository.Commit(tx)
		}
	}()

	amenities := make([]am.Amenity, 0)

	for _, am := range toCreate.Amenities {
		
		append(amenities, am.Amenity{
			Name: am,
			HotelId: hotel_id,
		})
	}

	room := mr.Room{
		Name: toCreate.Name,
		Price: new(big.Rat).SetString(toCreate.Price),
		HotelID: hotel_id,
	}

	if err := s.roomRepository.AddRoom(tx, room); err!=nil {
		return nil, err
	}
}

func (s *RoomServiceImpl)	UpdateRoom(toUpdate *dto.RoomRequestDTO, room_id int64) (*dto.RoomResponseDTO, error) {

}

