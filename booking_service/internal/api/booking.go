package api

import (
	"booking-service/internal/service"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	bookingService := service.NewBookingService(db)

	r.Get("/booking", func(w http.ResponseWriter, r *http.Request) {
		bookings, err := bookingService.GetAllBookings()
		if err != nil {
			http.Error(w, "Error retrieving bookings", http.StatusInternalServerError)
			return
		}
		jsonResponse(w, bookings)
	})

	return r
}

func jsonResponse(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
