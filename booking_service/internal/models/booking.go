package models

import "time"

type Booking struct {
	ID            int64     `json:"booking_id" gorm:"primaryKey"`
	RoomID        int64     `json:"room_id" gorm:"not null"`
	TimeFrom      time.Time `json:"time_from" gorm:"type:timestamptz"`
	TimeTo        time.Time `json:"time_to" gorm:"type:timestamptz"`
	ClientNumber  string    `json:"client_number" gorm:"not null"`
	TGUsername    string    `json:"tg_username" gorm:"not null"`
	BookingStatus string    `json:"status" gorm:"not null;default:'CREATED'"`
}
