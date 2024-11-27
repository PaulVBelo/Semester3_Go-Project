package api

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewRouter(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Get("/booking", func(w http.ResponseWriter, r *http.Request) {
		getAllBookings(w, r, db)
	})
	r.Get("/booking/{id}", func(w http.ResponseWriter, r *http.Request) {
		getBookingByID(w, r, db)
	})
	r.Get("/booking/room/{room_id}", func(w http.ResponseWriter, r *http.Request) {
		getCurrentForRoom(w, r, db)
	})
	r.Post("/booking", func(w http.ResponseWriter, r *http.Request) {
		createBooking(w, r, db)
	})
	r.Put("/booking/{id}", func(w http.ResponseWriter, r *http.Request) {
		updateBooking(w, r, db)
	})

	return r
}
