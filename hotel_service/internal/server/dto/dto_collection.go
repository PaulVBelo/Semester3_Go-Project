package dto

type RoomRequestDTO struct {
	Name		string					`json:"name,omitempty"`
	Price		string					`json:"price,omitempty"`
	Amenities	[]string				`json:"amenities,omitempty"`
}

type HotelRequestDTO struct {
	Name		string					`json:"name,omitempty"`
	Adress		string					`json:"adress,omitempty"`
	Rooms 		[]RoomRequestDTO		`json:"rooms,omitempty"`
}

type RoomResponseDTO struct {
	ID 			int64					`json:"id"`
	Name		string					`json:"name"`
	Price		string					`json:"price"`
	Amenities	[]string	`json:"amenities,omitempty"`
}

type HotelResponseDTO struct {
	ID			int64					`json:"id"`
	Name		string					`json:"name"`
	Adress		string					`json:"adress"`
	Rooms		[]RoomResponseDTO		`json:"rooms,omitempty"`
}