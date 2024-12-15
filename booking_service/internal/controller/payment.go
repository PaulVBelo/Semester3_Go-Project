package controller

import (
	"booking-service/internal/dto"
	"booking-service/internal/models"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"math"
	"net/http"
	"strconv"
	"time"
)

type FullRoomData struct {
	ID          int64    `json:"hotel_id"`
	Name        string   `json:"hotel_name"`
	Adress      string   `json:"hotel_adress"`
	PhoneNumber string   `json:"hotelier_number"`
	RoomId      int64    `json:"room_id"`
	RoomName    string   `json:"room_name"`
	Price       string   `json:"price"`
	Amenities   []string `json:"amenities,omitempty"`
}

func fetchRoomData(roomID int64) (*FullRoomData, error) {
	url := fmt.Sprintf("http://localhost:8080/api/book/%d", roomID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: received status code %d", resp.StatusCode)
	}

	var roomData FullRoomData
	if err := json.NewDecoder(resp.Body).Decode(&roomData); err != nil {
		return nil, err
	}

	return &roomData, nil
}

// PutPay обращается к hotel_service для получения данных комнаты и рассчитывает оплату
func (c *Controller) PutPay(w http.ResponseWriter, r *http.Request) {
	// Получение ID бронирования из URL
	bookingID := chi.URLParam(r, "id")
	if bookingID == "" {
		http.Error(w, "Booking ID is required", http.StatusBadRequest)
		return
	}

	// Получаем бронирование по ID из базы данных
	var booking models.Booking
	if err := c.DB.First(&booking, "id = ?", bookingID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "Booking not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Error finding booking: %v", err), http.StatusInternalServerError)
		return
	}

	// Обращение к API /api/book/:id для получения данных комнаты
	roomDTO, err := fetchRoomData(booking.RoomID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching room data: %v", err), http.StatusInternalServerError)
		return
	}

	// Расчет количества дней между time_from и time_to
	days := math.Ceil(float64(booking.TimeTo.Sub(booking.TimeFrom).Hours()) / 24)
	if days <= 0 {
		http.Error(w, "Invalid booking dates", http.StatusBadRequest)
		return
	}

	// Вычисляем стоимость проживания
	DayCost, err := strconv.Atoi(roomDTO.Price)

	// Возвращаем успешный ответ
	response := map[string]interface{}{
		"booking_id": booking.ID,
		"room_id":    booking.RoomID,
		"days":       days,
		"total_cost": float64(DayCost) * days,
		"status":     "ok",
	}

	var booking_data = dto.BookingData{
		ID:             roomDTO.ID,
		Name:           roomDTO.Name,
		Address:        roomDTO.Adress,
		HotelierNumber: roomDTO.PhoneNumber,
		RoomId:         booking.RoomID,
		RoomName:       roomDTO.Name,
		Payment:        strconv.Itoa(int(float64(DayCost) * days)),
		Amenities:      roomDTO.Amenities,
	}

	// Kafka event
	event := dto.BookingEventDTO{
		BookingId:    booking.ID,
		ClientNumber: booking.ClientNumber,
		TGUsername:   booking.TGUsername,
		TimeFrom:     booking.TimeFrom.Format(time.RFC3339),
		TimeTo:       booking.TimeTo.Format(time.RFC3339),
		Data:         booking_data,
	}

	fmt.Println(booking_data.Payment)

	if err := c.producer.Send(event); err != nil {
		http.Error(w, "Error sending Kafka event", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}