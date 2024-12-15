package controller

import (
	"booking-service/internal/models"
	"gorm.io/gorm"
	"time"
)

// CheckDateOverlap проверяет, есть ли пересечение дат бронирования для бронирования
func CheckDateOverlap(db *gorm.DB, roomID int64, timeFrom, timeTo time.Time) (bool, error) {
	var count int64

	err := db.Model(&models.Booking{}).
		Where("room_id = ? AND booking_status != ? AND booking_status != ? AND ((time_from < ? AND time_to > ?)"+
			" OR (time_from < ? AND time_to > ?))",
			roomID, "Сancelled", "Completed", timeTo, timeFrom, timeFrom, timeTo).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil // Если count > 0, то есть пересечение
}
