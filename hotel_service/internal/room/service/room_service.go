package service

import "hotel_service/internal/server/dto"

type RoomService interface {
	GetByID(id int64) (*dto.RoomResponseDTO, error)
	CreateRoom(toCreate *dto.RoomCreateRequestDTO, hotel_id int64) (*dto.RoomResponseDTO, error)
	UpdateRoom(toUpdate *dto.RoomUpdateRequestDTO, room_id int64) (*dto.RoomResponseDTO, error)
}
