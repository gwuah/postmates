package models

import "github.com/jinzhu/gorm"

type Electron struct {
	gorm.Model
	FirstName  string      `json:"firstName"`
	LastName   string      `json:"lastName"`
	MiddleName string      `json:"middleName"`
	Location   GeoLocation `json:"location"`
	Orders     []Order     `json:"orders"`
}
