package shared

import (
	"github.com/gwuah/api/database/models"
	"github.com/uber/h3-go"
)

type Meta struct {
	Type string `json:"type"`
}

type Coord struct {
	Longitude float64 `json:"longitude" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
}

type User struct {
	Id             string     `json:"id"`
	LastKnownIndex h3.H3Index `json:"lastKnownIndex"`
	Coord
}

type UserLocationUpdate struct {
	Id string `json:"id"`
	Coord
}

type DeliveryRequest struct {
	Meta        Meta   `json:"meta"`
	Origin      Coord  `json:"origin"`
	Destination Coord  `json:"destination"`
	ProductId   uint   `json:"productId"`
	Notes       string `json:"notes"`
	CustomerID  uint   `json:"customerId"`
}

type CancelDeliveryRequest struct {
	Meta   Meta `json:"meta"`
	TripId uint `json:"tripId"`
}

type GetClosestElectronsRequest struct {
	Meta   Meta   `json:"meta"`
	Id     string `json:"id"`
	Origin Coord  `json:"origin"`
}

type NewOrder struct {
	Meta  Meta          `json:"meta"`
	Order *models.Order `json:"order"`
}

type AcceptOrder struct {
	Meta    Meta `json:"meta"`
	OrderId uint `json:"orderId"`
}

type ElectronWithEta struct {
	Electron *User
	Duration float64
}
