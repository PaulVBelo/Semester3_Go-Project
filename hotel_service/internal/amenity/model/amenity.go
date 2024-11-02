package model

type Amenity struct {
	ID      int64  `gorm:"column:amenity_id;primaryKey"`
	Name    string `gorm:"column:ameity_name;size:128;not_null"`
	HotelID int64  `gorm:"column:hotel_id;"`
}
