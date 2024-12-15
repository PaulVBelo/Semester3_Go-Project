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

	//msg := tgbotapi.NewMessage(960397857, message)

	updates, err := bot.GetUpdates(tgbotapi.NewUpdate(0))
	if err != nil {
		logger.WithFields(logrus.Fields{
			"service": "delivery_system",
		}).Fatal("Error while getting updates:", err)
	}

	for _, update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID

			msg := tgbotapi.NewMessage(chatID, message)
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
	}

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
		bold := "\033[1m"
		reset := "\033[0m"
		message := bold + "📅 Booking Confirmation\n" + reset +
			"🏨 Hotel name: " + event.BookingData.HotelName + "\n" +
			"🛏️ Room name: " + event.BookingData.RoomName + "\n" +
			"⏰ Time from: " + event.TimeFrom +
			"⏰ Time to: " + event.TimeTo +
			"🔑 Booking ID: " + strconv.FormatInt(event.BookingId, 10) + "\n"

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
