package server

import (
	serviceH "hotel_service/internal/hotel/service"
	serviceR "hotel_service/internal/room/service"
	"hotel_service/internal/server/dto"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	router       *gin.Engine
	roomService  serviceR.RoomService
	hotelService serviceH.HotelService
}

func NewServer(hs serviceH.HotelService, rs serviceR.RoomService) *Server {
	router := gin.Default()
	s := &Server{
		router:       router,
		roomService:  rs,
		hotelService: hs,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.GET("api/hotels", s.getAll)
	s.router.GET("api/hotels/:id", s.getHotelByID)
	s.router.POST("api/hotels", s.createHotel)
	s.router.PUT("api/hotels/:id", s.updateHotel)
	s.router.GET("/rooms/:id", s.getRoomByID)
	s.router.POST("api/hotels/:id/room", s.createRoom)
	s.router.PUT("api/rooms/:id", s.updateRoom)
}

func (s *Server) getAll(c *gin.Context) {
	hotels, err := s.hotelService.GetAll()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not retrieved"})
		return
	}
	c.JSON(http.StatusOK, hotels)
}

func (s *Server) getHotelByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	hotel, err := s.hotelService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		return
	}
	c.JSON(http.StatusOK, hotel)
}

func (s *Server) createHotel(c *gin.Context) {
	var hotel dto.HotelCreateRequestDTO
	if err := c.ShouldBindJSON(&hotel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resp, err := s.hotelService.CreateHotel(&hotel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hotel", "details": err.Error()})
	}

	c.JSON(http.StatusCreated, resp)
}

func (s *Server) updateHotel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var hotel dto.HotelCreateRequestDTO
	if err := c.ShouldBindJSON(&hotel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resp, err := s.hotelService.UpdateHotel(id, &hotel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hotel", "details": err.Error()})
	}

	c.JSON(http.StatusCreated, resp)
}

func (s *Server) getRoomByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"id":    idStr,
			"error": err.Error(),
		}).Error("Invalid ID")

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	room, err := s.roomService.GetByID(id)
	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to find room")

		c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		return
	}

	logrus.WithTime(time.Now()).Info("Room successfully found")

	c.JSON(http.StatusOK, room)
}

func (s *Server) createRoom(c *gin.Context) {
	var room dto.RoomCreateRequestDTO
	if err := c.ShouldBindJSON(&room); err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Invalid Input")

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"id":    idStr,
			"error": err.Error(),
		}).Error("Invalid ID")

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	roomRsp, err := s.roomService.CreateRoom(&room, id)
	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to create room")

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room", "details": err.Error()})
	}

	logrus.WithTime(time.Now()).Info("Room successfully created")

	c.JSON(http.StatusCreated, roomRsp)
}

func (s *Server) updateRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"id":    idStr,
			"error": err.Error(),
		}).Error("Invalid ID")

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var room dto.RoomUpdateRequestDTO
	if err := c.ShouldBindJSON(&room); err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Invalid Input")

		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	roomRsp, err := s.roomService.UpdateRoom(&room, id)
	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to update room")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room", "details": err.Error()})
	}

	logrus.WithTime(time.Now()).Info("Room successfully updated")

	c.JSON(http.StatusOK, roomRsp)
}
