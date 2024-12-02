package server

import (
	serviceH "hotel_service/internal/hotel/service"
	serviceR "hotel_service/internal/room/service"
	"hotel_service/internal/server/dto"
	se "hotel_service/internal/server/errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	router       *gin.Engine
	roomService  serviceR.RoomService
	hotelService serviceH.HotelService
	responseTime *prometheus.HistogramVec
}

func NewServer(hs serviceH.HotelService, rs serviceR.RoomService) *Server {
	router := gin.Default()
	s := &Server{
		router:       router,
		roomService:  rs,
		hotelService: hs,
	}
	s.responseTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "hotel_svc_response_time",
		Help: "Successful HTTP response time for hotel svc",
		Buckets: []float64{0.1, 0.5, 1, 5, 10},
	}, []string{"method", "path"})
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.GET("api/hotels", s.getAll)
	s.router.GET("api/hotels/:id", s.getHotelByID)
	s.router.POST("api/hotels", s.createHotel)
	s.router.PUT("api/hotels/:id", s.updateHotel)
	s.router.GET("api/rooms/:id", s.getRoomByID)
	s.router.POST("api/hotels/:id/room", s.createRoom)
	s.router.PUT("api/rooms/:id", s.updateRoom)
	s.router.GET("api/book/:id", s.book)
	
	s.router.GET("/metrics", s.metrics)
}

func (s *Server) Run(port string) error {
	return s.router.Run(":" + port)
}

func (s *Server) getAll(c *gin.Context) {
	start := time.Now()

	hotels, err := s.hotelService.GetAll()
	if err != nil {
		switch err.(type) {
		case *se.NotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": "No hotels found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve hotelS",
				"details": "An unexpected error occurred"})
		}
		return
	}
	c.JSON(http.StatusOK, hotels)

	s.responseTime.WithLabelValues("GET", "/api/hotels").Observe(float64(time.Since(start).Seconds()))
}

func (s *Server) getHotelByID(c *gin.Context) {
	start := time.Now()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	hotel, err := s.hotelService.GetByID(id)
	if err != nil {
		switch err.(type) {
		case *se.NotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": "Hotel not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve hotel",
				"details": "An unexpected error occurred"})
		}
		return
	}
	c.JSON(http.StatusOK, hotel)

	s.responseTime.WithLabelValues("GET", "/api/hotels/:id").Observe(float64(time.Since(start).Seconds()))
}

func (s *Server) createHotel(c *gin.Context) {
	start := time.Now()

	var hotel dto.HotelCreateRequestDTO
	if err := c.ShouldBindJSON(&hotel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resp, err := s.hotelService.CreateHotel(&hotel)
	if err != nil {
		switch err.(type) {
		case *se.BadRequestError:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create hotel", "details": err.Error()})
		case *se.NotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to create hotel", "details": err.Error()})
		case *se.InternalServerError:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hotel"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create hotel",
				"details": "An unexpected error occurred"})
		}
		return
	}

	c.JSON(http.StatusCreated, resp)

	s.responseTime.WithLabelValues("POST", "/api/hotels").Observe(float64(time.Since(start).Seconds()))
}

func (s *Server) updateHotel(c *gin.Context) {
	start := time.Now()

	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var hotel dto.HotelUpdateRequestDTO
	if err := c.ShouldBindJSON(&hotel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	resp, err := s.hotelService.UpdateHotel(id, &hotel)
	if err != nil {
		switch err.(type) {
		case *se.BadRequestError:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update hotel", "details": err.Error()})
		case *se.InternalServerError:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hotel"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update hotel",
				"details": "An unexpected error occurred"})
		}
		return
	}

	c.JSON(http.StatusOK, resp)

	s.responseTime.WithLabelValues("PUT", "/api/hotels/:id").Observe(float64(time.Since(start).Seconds()))
}

func (s *Server) getRoomByID(c *gin.Context) {
	start := time.Now()

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

		switch err.(type) {
		case *se.NotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve room",
				"details": "An unexpected error occurred"})
		}
		return
	}

	logrus.WithTime(time.Now()).Info("Room successfully found")

	c.JSON(http.StatusOK, room)

	s.responseTime.WithLabelValues("GET", "/api/rooms/:id").Observe(float64(time.Since(start).Seconds()))
}

func (s *Server) createRoom(c *gin.Context) {
	start := time.Now()

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

		switch err.(type) {
		case *se.BadRequestError:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create room", "details": err.Error()})
		case *se.NotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to create room", "details": err.Error()})
		case *se.InternalServerError:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create room",
				"details": "An unexpected error occurred"})
		}
		return
	}

	logrus.WithTime(time.Now()).Info("Room successfully created")

	c.JSON(http.StatusCreated, roomRsp)

	s.responseTime.WithLabelValues("POST", "/api/hotels/ïd/room").Observe(float64(time.Since(start).Seconds()))
}

func (s *Server) updateRoom(c *gin.Context) {
	start := time.Now()

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

		switch err.(type) {
		case *se.BadRequestError:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to update room", "details": err.Error()})
		case *se.NotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": "Failed to update room", "details": err.Error()})
		case *se.InternalServerError:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update room",
				"details": "An unexpected error occurred"})
		}
		return
	}

	logrus.WithTime(time.Now()).Info("Room successfully updated")

	c.JSON(http.StatusOK, roomRsp)

	s.responseTime.WithLabelValues("PUT", "/api/rooms/:id").Observe(float64(time.Since(start).Seconds()))
}

func (s *Server) book(c *gin.Context) {
	start := time.Now()

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

	responseDTO, err := s.hotelService.GetExpandedRoomData(id)

	if err != nil {
		logrus.WithTime(time.Now()).WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Failed to find room")
		switch err.(type) {
		case *se.NotFoundError:
			c.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve room",
				"details": "An unexpected error occurred"})
		}
		return
	}

	logrus.WithTime(time.Now()).Info("Room successfully found")

	c.JSON(http.StatusOK, responseDTO)

	s.responseTime.WithLabelValues("GET", "/api/book/:id").Observe(float64(time.Since(start).Seconds()))
}


func (s *Server) metrics(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}