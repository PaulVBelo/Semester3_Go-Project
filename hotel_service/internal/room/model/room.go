package model

import (
  "database/sql/driver"
  "errors"
  "hotel_service/internal/amenity/model"
  "math/big"
  "gorm.io/gorm"
)

type Room struct {
  ID        int64                `gorm:"column:room_id;primaryKey"`
  Name      string               `gorm:"column:room_name;size:128;not null"`
  HotelID   int64                `gorm:"column:hotel_id;"`
  Price     BigRat               `gorm:"column:price;not null"`
  Amenities []*model.Amenity     `gorm:"many2many:room_x_amenity"`
}

type BigRat struct {
  Rat *big.Rat
}

func (br *BigRat) SetString(s string) (*BigRat, bool) {
  _, ok := br.Rat.SetString(s);
  return br, ok
}

func (br *BigRat) String() string {
  return br.Rat.String()
}

func (br *BigRat) Get() *big.Rat {
  return br.Rat
}

func (br *BigRat) Scan(value interface{}) error {
  if value == nil {
    br.Rat = big.NewRat(0, 1) 
    return nil
  }
  
  v, ok := value.(string)
  if !ok {
    return errors.New("failed to scan BigRat")
  }

  rat := new(big.Rat)
  if _, ok := rat.SetString(v); !ok {
    return errors.New("failed to parse BigRat from string")
  }
  br.Rat = rat
  return nil
}

func (br BigRat) Value() (driver.Value, error) {
  if br.Rat == nil {
    return nil, nil
  }
  return br.Rat.String(), nil
}

func (r *Room) BeforeSave(tx *gorm.DB) (err error) {
  if r.Price.Rat.Cmp(new(big.Rat)) != 1 {
    return errors.New("negative price not allowed")
  }
  return nil
}