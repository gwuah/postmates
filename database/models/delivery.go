package models

type State string

const (
	// delivery state types
	Pending         State = "pending"
	PendingPickup         = "pending_pickup"
	NearingPickup         = "nearing_pickup"
	AtPickup              = "at_pickup"
	DeliveryOngoing       = "delivery_ongoing"
	NearingDropoff        = "nearing_dropoff"
	AtDropOff             = "at_dropoff"
	Delivered             = "delivered"
	Cancelled             = "cancelled"

	// courier state types
	AwaitingDispatch = "awaiting_dispatch"
	Dispatched       = "dispatched"
	OnTrip           = "on_trip"
	Offline          = "offline"

	// vehicle state
	Inactive = "inactive"

	// customer state
	Searching = "searching"
	AtRest    = "atRest"
)

type Delivery struct {
	Model
	State                 State        `json:"state"`
	OriginLongitude       float64      `json:"originLongitude"`
	OriginLatitude        float64      `json:"originLatitude"`
	DestinationLongitude  float64      `json:"destinationLongitude"`
	DestinationLatitude   float64      `json:"destinationLatitude"`
	FinalCost             float64      `json:"finalCost"`
	InitialCost           float64      `json:"initialCost"`
	Completed             bool         `json:"completed"`
	Notes                 string       `json:"notes"`
	CustomerID            uint         `json:"customerId"`
	Customer              Customer     `json:"customer"`
	ProductID             uint         `json:"productId"`
	Product               Product      `json:"product"`
	CourierID             *uint        `json:"courierId,omitempty"`
	Courier               *Courier     `json:"courier,omitempty"`
	TripPoints            []*TripPoint `json:"tripPoints,omitempty"`
	CourierRating         float64      `json:"courierRating"`
	CourierRatingMessage  string       `json:"courierRatingMessage"`
	CustomerRating        float64      `json:"customerRating"`
	CustomerRatingMessage string       `json:"customerRatingMessage"`
}
