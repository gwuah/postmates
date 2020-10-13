package models

type Order struct {
	Model
	// Deliveries []Delivery `json:"deliveries"`
	CourierID uint    `json:"courierId"`
	Courier   Courier `json:"courier"`
	Completed bool    `json:"completed"`
	State     State   `json:"state"`
}
