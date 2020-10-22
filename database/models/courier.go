package models

type Courier struct {
	Model
	FirstName  string   `json:"firstName"`
	LastName   string   `json:"lastName"`
	MiddleName string   `json:"middleName"`
	Longitude  float64  `json:"longitude"`
	Latitude   float64  `json:"latitude"`
	State      State    `json:"state"`
	Vehicle    *Vehicle `json:"vehicle,omitempty"`
	// Deliveries []Delivery `json:"deliveries"`
	PhotoUrl string `json:"photoUrl"`
	Rating   int    `json:"rating"`
}
