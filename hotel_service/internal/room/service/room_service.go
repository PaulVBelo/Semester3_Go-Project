package service

import "hotel_service/internal/server/dto"

type RoomService interface {
	GetByID(id int64) (dto.RoomDTO, error)
	CreateRoom(toCreate *dto.RoomDTO) error
	UpdateRoom(toUpdate *dto.RoomDTO) error
}