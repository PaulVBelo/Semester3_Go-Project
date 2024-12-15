package models

import "time"

type Booking struct {
	ID            uint      `json:"booking_id" gorm:"primaryKey"`
	RoomID        uint      `json:"room_id" gorm:"not null"`
	TimeFrom      time.Time `json:"time_from" gorm:"type:timestamptz"`
	TimeTo        time.Time `json:"time_to" gorm:"type:timestamptz"`
	ClientNumber  string    `json:"client_number" gorm:"not null"`
	BookingStatus string    `json:"status" gorm:"not null;default:'CREATED'"`
}
