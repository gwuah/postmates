package models

import "github.com/jinzhu/gorm"

type Delivery struct {
	gorm.Model
	Status      string      `json:"status"`
	Origin      GeoLocation `json:"origin"`
	Destination GeoLocation `json:"destination"`
	Rating      float64     `json:"rating"`
	FinalCost   float64     `json:"finalCost"`
	InitialCost float64     `json:"initialCost"`
	Completed   bool        `json:"completed"`
	Notes       string      `json:"notes"`
	CustomerID  int         `json:"customerId"`
	Customer    Customer    `json:"customer"`
	OrderId     int         `json:"order_id"`
	Order       Order       `json:"order"`
}
