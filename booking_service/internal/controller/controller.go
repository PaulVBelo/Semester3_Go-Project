package controller

import (
	"booking-service/internal/config"
	"booking-service/internal/models"
	"booking-service/internal/utils"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type Controller struct {
	Config *config.Config
	DB     *gorm.DB
}

func NewController(cfg *config.Config, db *gorm.DB) *Controller {
	return &Controller{
		Config: cfg,
		DB:     db,
	}
}

// GetAllBookings возвращает все бронирования
func (c *Controller) GetAllBookings(w http.ResponseWriter, r *http.Request) {
	var bookings []models.Booking
	if err := c.DB.Find(&bookings).Error; err != nil {
		utils.ErrorResponse(w, "Error fetching bookings", err)
		return
	}
	utils.JSONResponse(w, bookings)
}

// GetBookingByID возвращает бронирование по ID
func (c *Controller) GetBookingByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var booking models.Booking
	if err := c.DB.First(&booking, id).Error; err != nil {
		utils.ErrorResponse(w, "Booking not found", err)
		return
	}
	utils.JSONResponse(w, booking)
}

// CreateBooking создает новое бронирование

func (ctrl *Controller) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var booking models.Booking
	if err := utils.ParseJSONBody(r, &booking); err != nil {
		utils.ErrorResponse(w, "Invalid request body", err)
		return
	}

	// Установка статуса по умолчанию, если не указан
	if booking.BookingStatus == "" {
		booking.BookingStatus = "CREATED"
	}

	// Сохранение бронирования в базе данных
	if err := ctrl.DB.Create(&booking).Error; err != nil {
		log.Printf("Error creating booking: %v", err)
		utils.ErrorResponse(w, "Error creating booking", err)
		return
	}

	// Ответ с успешно созданным бронированием
	utils.JSONResponse(w, booking)
}

// UpdateBooking обновляет существующее бронирование
func (c *Controller) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var booking models.Booking
	if err := c.DB.First(&booking, id).Error; err != nil {
		utils.ErrorResponse(w, "Booking not found", err)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&booking); err != nil {
		utils.ErrorResponse(w, "Invalid request body", err)
		return
	}

	if err := c.DB.Save(&booking).Error; err != nil {
		utils.ErrorResponse(w, "Error updating booking", err)
		return
	}
	utils.JSONResponse(w, booking)
}

// DeleteBooking удаляет бронирование по ID
func (c *Controller) DeleteBooking(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := c.DB.Delete(&models.Booking{}, id).Error; err != nil {
		utils.ErrorResponse(w, "Error deleting booking", err)
		return
	}
	utils.JSONResponse(w, map[string]string{"message": "Booking deleted"})
}
