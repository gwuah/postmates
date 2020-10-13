package shared

import (
	"errors"

	"github.com/gwuah/api/database/models"
	"github.com/uber/h3-go"
)

var (
	MAPBOX_ERROR = errors.New("Mapbox Request failed")
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
	Id         string       `json:"id"`
	State      models.State `json:"state"`
	DeliveryId uint         `json:"deliveryId"`
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

type GetClosestCouriersRequest struct {
	Origin Coord `json:"origin"`
}

type NewDelivery struct {
	Meta             Meta             `json:"meta"`
	Delivery         *models.Delivery `json:"delivery"`
	DistanceToPickup float64          `json:"distanceToPickup"`
	DurationToPickup float64          `json:"durationToPickup"`
}

type AcceptDelivery struct {
	Meta       Meta `json:"meta"`
	DeliveryId uint `json:"deliveryId"`
}

type CourierWithEta struct {
	Courier  *User   `json:"courier"`
	Distance float64 `json:"distance"`
	Duration float64 `json:"duration"`
}

type DeliveryAcceptedPayload struct {
	Meta             Meta            `json:"meta"`
	Courier          models.Courier  `json:"courier"`
	Delivery         models.Delivery `json:"delivery"`
	DistanceToPickup float64         `json:"distanceToPickup"`
	DurationToPickup float64         `json:"durationToPickup"`
}

type NoCourierAvailable struct {
	Meta    Meta   `json:"meta"`
	Message string `json:"message"`
}

type CourierLocation struct {
	Meta Meta `json:"meta"`
	Coord
	DistanceToPickup float64 `json:"distanceToPickup"`
	DurationToPickup float64 `json:"durationToPickup"`
}

type PricePerProduct struct {
	ProductId uint `json:"productId"`
	Price     int  `json:"price"`
}

type GetDeliveryCostRequest struct {
	Origin      Coord `json:"origin" validate:"required"`
	Destination Coord `json:"destination" validate:"required"`
}
