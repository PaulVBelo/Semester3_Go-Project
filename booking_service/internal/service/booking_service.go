package service

import (
	"booking-service/internal/models"
	"database/sql"
)

type BookingService struct {
	DB *sql.DB
}

func NewBookingService(db *sql.DB) *BookingService {
	return &BookingService{DB: db}
}

func (s *BookingService) GetAllBookings() ([]models.Booking, error) {
	rows, err := s.DB.Query("SELECT * FROM booking")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookings []models.Booking
	for rows.Next() {
		var b models.Booking
		if err := rows.Scan(&b.ID, &b.RoomID, &b.TimeFrom, &b.TimeTo, &b.ClientNumber, &b.BookingStatus); err != nil {
			return nil, err
		}
		bookings = append(bookings, b)
	}
	return bookings, nil
}
