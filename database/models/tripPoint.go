package models

type TripPoint struct {
	Model
	DeliveryID uint      `json:"deliveryId"`
	Delivery   *Delivery `json:"delivery"`
	Longitude  float64   `json:"longitude"`
	Latitude   float64   `json:"latitude"`
	State      State     `json:"state"`
}
