package model

import (
	"database/sql"
	"fmt"
	"time"
)

type Booking struct {
	BookingID     int64     `json:"booking_id"`
	RoomID        int64     `json:"room_id"`
	TimeFrom      time.Time `json:"time_from"`
	TimeTo        time.Time `json:"time_to"`
	ClientNumber  string    `json:"client_number"`
	BookingStatus string    `json:"booking_status"`
}

func GetAllBookings(db *sql.DB) ([]Booking, error) {
	rows, err := db.Query("SELECT * FROM booking")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var booking Booking
		if err := rows.Scan(&booking.BookingID, &booking.RoomID, &booking.TimeFrom, &booking.TimeTo, &booking.ClientNumber, &booking.BookingStatus); err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func GetBookingByID(db *sql.DB, id string) (*Booking, error) {
	var booking Booking
	err := db.QueryRow("SELECT * FROM booking WHERE booking_id = $1", id).
		Scan(&booking.BookingID, &booking.RoomID, &booking.TimeFrom, &booking.TimeTo, &booking.ClientNumber, &booking.BookingStatus)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("no booking found with id: %s", id)
	}
	return &booking, err
}

func GetCurrentBookingsForRoom(db *sql.DB, roomID string) ([]Booking, error) {
	now := time.Now()
	rows, err := db.Query(
		"SELECT * FROM booking WHERE room_id = $1 AND time_to > $2 AND booking_status != 'CANCELLED'",
		roomID, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []Booking
	for rows.Next() {
		var booking Booking
		if err := rows.Scan(&booking.BookingID, &booking.RoomID, &booking.TimeFrom, &booking.TimeTo, &booking.ClientNumber, &booking.BookingStatus); err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	return bookings, nil
}

func CreateBooking(db *sql.DB, booking *Booking) (int64, error) {
	query := `
		INSERT INTO booking (room_id, time_from, time_to, client_number, booking_status)
		VALUES ($1, $2, $3, $4, $5) RETURNING booking_id
	`
	var bookingID int64
	err := db.QueryRow(query, booking.RoomID, booking.TimeFrom, booking.TimeTo, booking.ClientNumber, booking.BookingStatus).Scan(&bookingID)
	return bookingID, err
}

func UpdateBookingStatus(db *sql.DB, id string, status string) error {
	_, err := db.Exec("UPDATE booking SET booking_status = $1 WHERE booking_id = $2", status, id)
	return err
}
