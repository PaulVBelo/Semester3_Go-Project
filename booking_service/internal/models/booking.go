package models

type Booking struct {
	ID            int    `json:"id"`
	RoomID        int    `json:"room_id"`
	TimeFrom      string `json:"time_from"` // ISO 8601 format
	TimeTo        string `json:"time_to"`   // ISO 8601 format
	ClientNumber  string `json:"client_number"`
	BookingStatus string `json:"booking_status"`
}
