package repository

import (
	"hotel_service/internal/room/model"

	"gorm.io/gorm"
)

type RoomRepository interface {
	Begin() (*gorm.DB, error)
	Commit(tx *gorm.DB) error
	Rollback(tx *gorm.DB) error

	AddRoom(tx *gorm.DB, room *model.Room) error
	AddRooms(tx *gorm.DB, rooms *[]model.Room) error
	GetRoomById(id int64) (*model.Room, error)
	UpdateRoom(tx *gorm.DB, room *model.Room) error
	GetAll() ([]model.Room, error)
}