package service

import "hotel_service/internal/server/dto"

type RoomService interface {
	GetByID(id int64) (*dto.RoomResponseDTO, error)
	CreateRoom(toCreate *dto.RoomRequestDTO, hotel_id int64) (*dto.RoomResponseDTO, error)
	UpdateRoom(toUpdate *dto.RoomRequestDTO, room_id int64) (*dto.RoomResponseDTO, error)
}