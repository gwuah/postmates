package models

import (
	"time"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
}

type Customer struct {
	Model
	State     State   `json:"state"`
	Phone     string  `gorm:"not null;unique" json:"phone"`
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Email     string  `json:"email"`
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
	Code      int     `json:"code"`
	Active    bool    `json:"active" gorm:"default=false"`
	Token     string  `json:"-"`
	Rating    int     `json:"rating"`
}
