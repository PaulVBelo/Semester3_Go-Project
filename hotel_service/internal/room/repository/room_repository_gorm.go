package repository

import (
	"gorm.io/gorm"
	"hotel_service/internal/room/model"
)

type roomRepositoryWithGorm struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &roomRepositoryWithGorm{db: db}
}

func (r *roomRepositoryWithGorm) AddRoom(room *model.Room) error {
	return r.db.Create(room).Error
}

func (r *roomRepositoryWithGorm) AddRooms(rooms *[]model.Room) error {
	return r.db.Create(rooms).Error
}

func (r *roomRepositoryWithGorm) GetRoomById(id int64) (*model.Room, error) {
	var room model.Room
	if err := r.db.Preload("Amenities").First(&room, id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepositoryWithGorm) UpdateRoom(room *model.Room) error {
	return r.db.Save(room).Error
}

func (r *roomRepositoryWithGorm) GetAll() ([]model.Room, error) {
	var rooms []model.Room
	result := r.db.Find(&rooms)
	return rooms, result.Error
}