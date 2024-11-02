package model

import (
	"errors"
	"hotel_service/internal/amenity/model"
	"math/big"
	"gorm.io/gorm"
)

type Room struct {
	ID 			int64 				`gorm:"column:room_id;primaryKey"`
	Name 		string				`gorm:"column:room_name;size:128;not_null"`
	HotelID		int64				`gorm:"column:hotel_id;"`
	Price		big.Rat				`gorm:"column:price;not_null"`
	
	Amenities	[]*model.Amenity	`gorm:"many2many:room_x_amenity"`
}

func(r *Room) BeforeSave(tx *gorm.DB) (err error) {
	if r.Price.Cmp(new(big.Rat)) != 1 {
		return errors.New("negative price not allowed")
	}
	return nil
}