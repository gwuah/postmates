package models

import "github.com/jinzhu/gorm"

type Customer struct {
	gorm.Model
	ID        uint    `gorm:"primary_key" json:"id"`
	Status    string  `json:"status"`
	Phone     string  `json:"phone"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Email     string  `json:"email"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}
