package models

type Delivery struct {
	Model
	Status               string   `json:"status"`
	OriginLongitude      float64  `json:"originLongitude"`
	OriginLatitude       float64  `json:"originLatitude"`
	DestinationLongitude float64  `json:"destinationLongitude"`
	DestinationLatitude  float64  `json:"destinationLatitude"`
	Rating               float64  `json:"rating"`
	FinalCost            float64  `json:"finalCost"`
	InitialCost          float64  `json:"initialCost"`
	Completed            bool     `json:"completed"`
	Notes                string   `json:"notes"`
	CustomerID           int      `json:"customerId"`
	Customer             Customer `json:"customer"`
	OrderID              int      `json:"orderId"`
	Order                Order    `json:"order"`
}
