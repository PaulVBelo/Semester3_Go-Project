package model

import (
	"hotel_service/internal/room/model"
	se "hotel_service/internal/server/errors"
	"regexp"

	"gorm.io/gorm"
)

type Hotel struct {
	ID 			int64 			`gorm:"column:hotel_id;primaryKey"`
	Name 		string			`gorm:"column:hotel_name;uniqueIndex:idx_name_adress;size:64;not_null"`
	Adress		string			`gorm:"column:hotel_adress;uniqueIndex:idx_name_adress;size:256;not_null"`
	PhoneNumber	string			`gorm:"column:phone_number;size:32;not_null"`

	Rooms		[]*model.Room	`gorm:"foreignKey:HotelID"`
}

func(h *Hotel) BeforeSave(tx *gorm.DB) (err error) {
	if !isValid(h.PhoneNumber) {
		return &se.BadRequestError{"Invalid phone number"}
	}
	return nil
}

func isValid(phone_number string) bool {
	re := regexp.MustCompile(`^\+?\d+$`)
	return re.MatchString(phone_number)
}