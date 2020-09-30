package models

type Order struct {
	Model
	Deliveries []Delivery `json:"deliveries"`
	ElectronID uint       `json:"electronId"`
	Electron   Electron   `json:"electron"`
	Completed  bool       `json:"completed"`
	Status     Status     `json:"status"`
}
