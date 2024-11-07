package dto

type RoomRequestDTO struct {
	Name		string
	Price		string
	Rooms		[]RoomRequestDTO
}

type HotelRequestDTO struct {
	Name		string
	Adress		string
}

type RoomResponseDTO struct {
	ID 			int64
	Name		string
	HotelID 	int64
	Price		string
	Amenities	[]AmenityResponseDTO	`json:"omitempty"`
}

type AmenityResponseDTO struct {
	ID 			int64
	Name		string
}

type HotelResponseDTO struct {
	ID			int64
	Name		string
	Adress		string
	Rooms		[]RoomResponseDTO
	Amenities	[]AmenityResponseDTO
}