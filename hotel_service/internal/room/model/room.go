package model

import (
	"errors"
	"gorm.io/gorm"
	"math/big"
)

type Room struct {
	ID 			int64 			`gorm:"column:room_id;primaryKey"`
	Name 		string			`gorm:"column:room_name;not_null"`
	HotelID		int64			`gorm:"column:hotel_id;"`
	Price		big.Rat			`gorm:"column:price;not_null"`
}

func(r *Room) BeforeSave(tx *gorm.DB) (err error) {
	if r.Price.Cmp(new(big.Rat)) != 1 {
		return errors.New("negative price not allowed")
	}
	return nil
}