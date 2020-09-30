package models

var (
	STATUS_TYPES = map[string]string{
		"pending":          "pending",
		"pending_pickup":   "pending_pickup",
		"nearing_pickup":   "nearing_pickup",
		"delivery_ongoing": "delivery_ongoing",
		"nearing_dropoff":  "nearing_dropoff",
		"delivered":        "delivered",
		"canceled":         "canceled",
	}
)

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
	CustomerID           uint     `json:"customerId"`
	Customer             Customer `json:"customer"`
	OrderID              uint     `json:"orderId"`
	Order                Order    `json:"order"`
	ProductID            uint     `json:"productId"`
	Product              Product  `json:"product"`
}
