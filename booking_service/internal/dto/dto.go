package dto

type BookingData struct {
	ID     			int64             		`json:"hotel_id"`
	Name   			string            		`json:"hotel_name"`
	Adress 			string            		`json:"hotel_adress"`
	HotelierNumber		string				`json:"hotelier_number"`
	RoomId			int64					`json:"room_id"`
	RoomName      	string   				`json:"room_name"`
	Payment    		string   				`json:"payment"`
	Amenities 		[]string 				`json:"amenities,omitempty"`
}

type BookingEventDTO struct {
	BookingId		int64					`json:"booking_id"`
	ClientNumber	string					`json:"client_phone_number"`
	TimeFrom		string					`json:"time_from"`
	TimeTo			string					`json:"time_to"`
	Data			BookingData				`json:"booking_data"`
}