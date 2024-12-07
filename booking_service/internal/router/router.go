package router

import (
	"booking-service/internal/controller"
	"booking-service/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func NewRouter(ctrl *controller.Controller) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Route("/booking", func(r chi.Router) {
		r.Get("/", ctrl.GetAllBookings)
		r.Get("/{id}", ctrl.GetBookingByID)
		r.Post("/", ctrl.CreateBooking)
		r.Put("/{id}", ctrl.UpdateBooking)
		r.Delete("/{id}", ctrl.DeleteBooking)
	})

	return router
}
