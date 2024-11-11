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
	Name   string                 `json:"name,omitempty"`
	Adress string                 `json:"adress,omitempty"`
	Rooms  []RoomCreateRequestDTO `json:"rooms,omitempty"`
}

type RoomResponseDTO struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	Price     string   `json:"price"`
	Amenities []string `json:"amenities,omitempty"`
}

type HotelResponseDTO struct {
	ID     int64             `json:"id"`
	Name   string            `json:"name"`
	Adress string            `json:"adress"`
	Rooms  []RoomResponseDTO `json:"rooms,omitempty"`
}
