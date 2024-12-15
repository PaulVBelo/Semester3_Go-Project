package controller

import (
	"booking-service/internal/config"
	"booking-service/internal/models"
	"booking-service/internal/producer"
	"booking-service/internal/utils"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
)

type Controller struct {
	Config   *config.Config
	DB       *gorm.DB
	producer *producer.BookingEventProducer
}

func NewController(cfg *config.Config, db *gorm.DB, producer *producer.BookingEventProducer) *Controller {
	return &Controller{
		Config:   cfg,
		DB:       db,
		producer: producer,
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
	}
	utils.JSONResponse(w, booking)
}

// CreateBooking создает новое бронирование
func (c *Controller) CreateBooking(w http.ResponseWriter, r *http.Request) {
	var booking models.Booking

	// Чтение JSON тела запроса
	err := json.NewDecoder(r.Body).Decode(&booking)
	if err != nil {
		utils.ErrorResponse(w, "Invalid request payload", err)
		return
	}

	// Проверка на пересечение дат
	overlap, err := CheckDateOverlap(c.DB, booking.RoomID, booking.TimeFrom, booking.TimeTo)
	if err != nil {
		utils.ErrorResponse(w, "Error checking date overlap", err)
		return
	}

	if overlap {
		utils.ErrorResponse(w, "Time overlap with another booking", nil)
		return
	}

	// Запись в базу данных
	booking.BookingStatus = "created" // Статус по умолчанию

	err = c.DB.Create(&booking).Error
	if err != nil {
		utils.ErrorResponse(w, "Error creating booking", err)
		return
	}

	// Ответ с успешным созданием
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

//// DeleteBooking удаляет бронирование по ID
//func (c *Controller) DeleteBooking(w http.ResponseWriter, r *http.Request) {
//	id := chi.URLParam(r, "id")
//	if err := c.DB.Delete(&models.Booking{}, id).Error; err != nil {
//		utils.ErrorResponse(w, "Error deleting booking", err)
//		return
//	}
//	utils.JSONResponse(w, map[string]string{"message": "Booking deleted"})
//}
