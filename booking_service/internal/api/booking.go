package api

import (
	"booking-service/internal/model"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func getAllBookings(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	bookings, err := model.GetAllBookings(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(bookings)
}

func getBookingByID(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := chi.URLParam(r, "id")
	booking, err := model.GetBookingByID(db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(booking)
}

func getCurrentForRoom(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	roomID := chi.URLParam(r, "room_id")
	bookings, err := model.GetCurrentBookingsForRoom(db, roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(bookings)
}

func createBooking(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var booking model.Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	booking.TimeFrom = time.Now() // Пример установки времени
	bookingID, err := model.CreateBooking(db, &booking)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"booking_id": bookingID})
}

func updateBooking(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	id := chi.URLParam(r, "id")
	var status struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := model.UpdateBookingStatus(db, id, status.Status); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
