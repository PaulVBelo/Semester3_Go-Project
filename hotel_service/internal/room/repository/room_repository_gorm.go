package repository

import (
	"hotel_service/internal/room/model"

	"gorm.io/gorm"
)

type RoomRepositoryWithGorm struct {
	db *gorm.DB
}

func NewRoomRepository(db *gorm.DB) RoomRepository {
	return &RoomRepositoryWithGorm{db: db}
}

func (r *RoomRepositoryWithGorm) Begin() (*gorm.DB, error) {
	return r.db.Begin(), nil
}

func (r *RoomRepositoryWithGorm) Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (r *RoomRepositoryWithGorm) Rollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func (r *RoomRepositoryWithGorm) AddRoom(tx *gorm.DB, room *model.Room) error {
	return tx.Create(room).Error
}

func (r *RoomRepositoryWithGorm) AddRooms(tx *gorm.DB, rooms *[]model.Room) error {
	return tx.Create(rooms).Error
}

func (r *RoomRepositoryWithGorm) GetRoomById(id int64) (*model.Room, error) {
	var room model.Room
	if err := r.db.Preload("Amenities").First(&room, id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RoomRepositoryWithGorm) UpdateRoom(tx *gorm.DB, room *model.Room) error {
	return tx.Save(room).Error
}

func (r *RoomRepositoryWithGorm) GetAll() ([]model.Room, error) {
	var rooms []model.Room
	result := r.db.Find(&rooms)
	return rooms, result.Error
}
