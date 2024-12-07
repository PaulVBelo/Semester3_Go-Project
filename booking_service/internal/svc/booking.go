package svc

import (
	"booking-service/internal/models"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
)

type BookingHandler struct {
	db *gorm.DB
}

func NewBookingHandler(db *gorm.DB) *BookingHandler {
	return &BookingHandler{db: db}
}

func (h *BookingHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	var bookings []models.Booking
	if result := h.db.Find(&bookings); result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, bookings)
}

func (h *BookingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var booking models.Booking
	if result := h.db.First(&booking, id); errors.Is(result.Error, gorm.ErrRecordNotFound) {
		respondWithError(w, http.StatusNotFound, "Booking not found")
		return
	} else if result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, booking)
}

func (h *BookingHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var booking models.Booking
	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	booking.BookingStatus = "CREATED"
	if result := h.db.Create(&booking); result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, booking)
}

func (h *BookingHandler) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var updateRequest struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateRequest); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	if result := h.db.Model(&models.Booking{}).Where("id = ?", id).Update("status", updateRequest.Status); result.Error != nil {
		respondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}
