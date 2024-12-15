package handler

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"

	"delivery_system/proto/gen"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendMessageToTelegram(username, message string) error {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	botToken := os.Getenv("API_TOKEN")
	if botToken == "" {
		logger.WithFields(logrus.Fields{
			"service": "delivery_system",
		}).Error("Telegram bot token is not set")
		return nil
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"service": "delivery_system",
			"error":   err,
		}).Error("Failed to create Telegram bot")
		return err
	}

	msg := tgbotapi.NewMessageToChannel("@"+username, message)

	_, err = bot.Send(msg)
	if err != nil {
		logger.WithFields(logrus.Fields{
			"service": "delivery_system",
			"error":   err,
		}).Error("Failed to send Telegram message")
		return err
	}

	logger.WithFields(logrus.Fields{
		"service":  "delivery_system",
		"username": username,
	}).Info("Telegram message sent successfully")
	return nil
}

func HandleBookingEvent(event *gen.BookingEvent) error {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.WithFields(logrus.Fields{
		"service": "delivery_system",
		"event":   event,
	}).Info("Received booking event")

	if event.TgUsername != "" {
		message := "Booking Confirmation\n" +
			"Hotel name: " + event.BookingData.HotelName + "\n" +
			"Room name: " + event.BookingData.RoomName + "\n" +
			"Booking ID: " + strconv.FormatInt(event.BookingId, 10) + "\n"

		err := sendMessageToTelegram(event.TgUsername, message)
		if err != nil {
			logger.WithFields(logrus.Fields{
				"service":  "delivery_system",
				"username": event.TgUsername,
				"error":    err,
			}).Error("Failed to send message to Telegram")
		}
	}

	return nil
}
