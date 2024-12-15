package main

import (
	modelA "hotel_service/internal/amenity/model"
	repositoryA "hotel_service/internal/amenity/repository"
	modelH "hotel_service/internal/hotel/model"
	repositoryH "hotel_service/internal/hotel/repository"
	svcH "hotel_service/internal/hotel/service"
	modelR "hotel_service/internal/room/model"
	repositoryR "hotel_service/internal/room/repository"
	svcR "hotel_service/internal/room/service"
	"hotel_service/internal/server"
	"os"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	if err := godotenv.Load("../.env.dev"); err != nil {
		logrus.Fatal("Error loading .env file")
	}
	
	dsn := "host=" + os.Getenv("DB_HOST") + 
		" port=" + os.Getenv("DB_PORT") + 
		" user=" + os.Getenv("DB_USER") + 
		" password=" + os.Getenv("DB_PASSWORD") + 
		" dbname=" + os.Getenv("DB_NAME") + " sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.WithFields(logrus.Fields{"error": err,}).Fatal("failed to connect to database")
	}

	if err := db.AutoMigrate(&modelH.Hotel{}, &modelR.Room{}, &modelA.Amenity{}); err != nil {
		logrus.WithFields(logrus.Fields{"error": err,}).Fatal("failed to migrate database")
	}

	roomRepo := repositoryR.NewRoomRepository(db)
	hotelRepo := repositoryH.NewHotelRepository(db)
	amenityRepo := repositoryA.NewAmenityRepository(db)
	
	roomSvc := svcR.NewRoomService(roomRepo, amenityRepo)
	hotelSvc := svcH.NewHotelService(roomRepo, amenityRepo, hotelRepo)

	server := server.NewServer(hotelSvc, roomSvc)
	port := os.Getenv("PORT")
	if err := server.Run(port); err != nil {
		logrus.WithFields(logrus.Fields{"error": err,}).Fatal("failed to run server")
	}
}