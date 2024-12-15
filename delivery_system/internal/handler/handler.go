package handler

import (
	"log"
	"os"
	"strconv"

	"delivery_system/proto/gen"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func sendMessageToTelegram(username, message string) error {
	botToken := os.Getenv("API_TOKEN")
	if botToken == "" {
		log.Println("Telegram bot token is not set")
		return nil
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Printf("Failed to create Telegram bot: %v", err)
		return err
	}

	msg := tgbotapi.NewMessageToChannel("@"+username, message)

	_, err = bot.Send(msg)
	if err != nil {
		log.Printf("Failed to send Telegram message: %v", err)
		return err
	}

	log.Printf("Telegram message sent successfully to user: @%s", username)
	return nil
}

func HandleBookingEvent(event *gen.BookingEvent) error {
	log.Printf("Received booking event: %+v\n", event)

	if event.TgUsername != "" {
		message := "Booking Confirmation\n" +
			"Hotel name: " + event.BookingData.HotelName + "\n" +
			"Room name: " + event.BookingData.RoomName + "\n" +
			"Booking ID: " + strconv.FormatInt(event.BookingId, 10) + "\n"

		err := sendMessageToTelegram(event.TgUsername, message)
		if err != nil {
			log.Printf("Failed to send message to Telegram: %v", err)
		}
	}

	return nil
}
