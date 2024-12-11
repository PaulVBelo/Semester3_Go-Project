package dto

type RoomCreateRequestDTO struct {
	Name      string   `json:"name"`
	Price     string   `json:"price"`
	Amenities []string `json:"amenities,omitempty"`
}

type RoomUpdateRequestDTO struct {
	Name      *string  `json:"name,omitempty"`
	Price     *string  `json:"price,omitempty"`
	Amenities []string `json:"amenities,omitempty"`
}

type HotelCreateRequestDTO struct {
	Name   	string                 	`json:"name"`
	Adress 	string                 	`json:"adress"`
	PhoneNumber	string				`json:"phone_number"`
	Rooms  	[]RoomCreateRequestDTO 	`json:"rooms,omitempty"`
}

type HotelUpdateRequestDTO struct {
	Name   		*string                 	`json:"name"`
	Adress 		*string                 	`json:"adress"`
	PhoneNumber	*string						`json:"phone_number"`
}

type RoomResponseDTO struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	Price     string   `json:"price"`
	Amenities []string `json:"amenities,omitempty"`
}

type HotelResponseDTO struct {
	ID     int64             		`json:"id"`
	Name   string            		`json:"name"`
	Adress string            		`json:"adress"`
	PhoneNumber	string				`json:"phone_number"`
	Rooms  []*RoomResponseDTO 		`json:"rooms,omitempty"`
}

type HotelShortResponseDTO struct {
	ID     int64             		`json:"id"`
	Name   string            		`json:"name"`
	Adress string            		`json:"adress"`
	PhoneNumber	string				`json:"phone_number"`
}

type FullRoomData struct {
	ID     			int64             		`json:"hotel_id"`
	Name   			string            		`json:"hotel_name"`
	Adress 			string            		`json:"hotel_adress"`
	PhoneNumber		string					`json:"hotelier_number"`
	RoomId			int64					`json:"room_id"`
	RoomName      	string   				`json:"room_name"`
	Price     		string   				`json:"price"`
	Amenities 		[]string 				`json:"amenities,omitempty"`
}