package shared

type Coord struct {
	Longitude float64 `json:"longitude" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
}

type DeliveryRequest struct {
	BaseMessage BaseMessage `json:"meta"`
	Origin      Coord       `json:"origin"`
	Destination Coord       `json:"destination"`
	ProductId   uint        `json:"productId"`
	Notes       string      `json:"notes"`
	CustomerID  uint        `json:"customerId"`
}

type CancelDeliveryRequest struct {
	BaseMessage BaseMessage `json:"meta"`
	TripId      uint        `json:"tripId"`
}

type BaseMessage struct {
	Type string `json:"type"`
}
