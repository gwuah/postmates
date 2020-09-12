package models

import "github.com/jinzhu/gorm"

type Order struct {
	gorm.Model
	Deliveries []Delivery `json:"deliveries"`
	ElectronID int        `json:"electronId"`
	Electron   Electron   `json:"electron"`
	Completed  bool       `json:"completed"`
}
