package server

import (
	serviceH "hotel_service/internal/hotel/service"
	serviceR "hotel_service/internal/room/service"
	"hotel_service/internal/server/dto"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	s.router.GET("/hotels/:id", s.getHotelByID)
	s.router.POST("/hotels", s.createHotel)
	s.router.PUT("/hotels/:id", s.updateHotel)
	s.router.GET("/rooms/:id", s.getRoomByID)
	s.router.POST("/hotels/:id/room", s.createRoom)
	s.router.PUT("/rooms/:id", s.updateRoom)
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
	var hotel dto.HotelRequestDTO
	if err := c.ShouldBindJSON(&hotel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := s.hotelService.CreateHotel(&hotel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tvoi soft gavno"})
	}

	c.JSON(http.StatusCreated, hotel)
}

func (s *Server) updateHotel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var hotel dto.HotelRequestDTO
	if err := c.ShouldBindJSON(&hotel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := s.hotelService.UpdateHotel(id, &hotel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tvoi soft gavno"})
	}

	c.JSON(http.StatusCreated, hotel)
}

func (s *Server) getRoomByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	room, err := s.roomService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		return
	}
	c.JSON(http.StatusOK, room)
}

func (s *Server) createRoom(c *gin.Context) {
	var room dto.RoomRequestDTO
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	room.HotelID = id

	if err := s.roomService.CreateRoom(&room); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tvoi soft gavno"})
	}

	c.JSON(http.StatusCreated, room)
}

func (s *Server) updateRoom(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var room dto.RoomRequestDTO
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	var roomRsp dto.RoomResponseDTO
	if roomRsp, err := s.roomService.UpdateRoom(&room, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Tvoi soft gavno"})
	}

	c.JSON(http.StatusCreated, room)
}
