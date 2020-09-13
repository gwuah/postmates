package models

type Order struct {
	Model
	Deliveries []Delivery `json:"deliveries"`
	ElectronID *int       `json:"electronId"`
	Electron   *Electron  `json:"electron"`
	Completed  bool       `json:"completed"`
}
