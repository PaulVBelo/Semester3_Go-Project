package repository

import "hotel_service/internal/room/model"

type RoomRepository interface {
	AddRoom(room *model.Room) error
	AddRooms(rooms *[]model.Room) error
	GetRoomById(id int64) (*model.Room, error)
	UpdateRoom(room *model.Room) error
	GetAll() ([]model.Room, error)
}