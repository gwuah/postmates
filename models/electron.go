package models

import "github.com/jinzhu/gorm"

type Electron struct {
	gorm.Model
	FirstName  string  `json:"firstName"`
	LastName   string  `json:"lastName"`
	MiddleName string  `json:"middleName"`
	Longitude  float64 `json:"longitude"`
	Latitude   float64 `json:"latitude"`
	Orders     []Order `json:"orders"`
}
