package dto

type RoomDTO struct {
	ID 			int64
	Name		string
	HotelID 	int64
	Price		string
	Amenities	[]AmenityDTO	`json:"omitempty"`
}

type AmenityDTO struct {
	ID 			int64
	Name		string
}

type HotelDTO struct {
	ID			int64
	Name		string
	Adress		string
	Rooms		[]RoomDTO
	Amenities	[]AmenityDTO
}