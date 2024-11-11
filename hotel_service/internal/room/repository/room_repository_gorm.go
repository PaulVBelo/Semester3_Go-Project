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

func (r *roomRepositoryWithGorm) Begin() (*gorm.DB, error) {
	return r.db.Begin(), nil
}

func (r *roomRepositoryWithGorm) Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}
	
func (r *roomRepositoryWithGorm) Rollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func ( r *roomRepositoryWithGorm) AddRoom(tx *gorm.DB, room *model.Room) error {
	return tx.Create(room).Error
}

func (r *roomRepositoryWithGorm) AddRooms(tx *gorm.DB, rooms *[]model.Room) error {
	return tx.Create(rooms).Error
}

func (r *roomRepositoryWithGorm) GetRoomById(id int64) (*model.Room, error) {
	var room model.Room
	if err := r.db.Preload("Amenities").First(&room, id).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *roomRepositoryWithGorm) UpdateRoom(tx *gorm.DB, room *model.Room) error {
	return tx.Save(room).Error
}

func (r *roomRepositoryWithGorm) GetAll() ([]model.Room, error) {
	var rooms []model.Room
	result := r.db.Find(&rooms)
	return rooms, result.Error
}