package models

type Status string

const (
	// delivery status types
	Pending         Status = "pending"
	PendingPickup          = "pending_pickup"
	NearingPickup          = "nearing_pickup"
	AtPickup               = "at_pickup"
	DeliveryOngoing        = "delivery_ongoing"
	NearingDropoff         = "nearing_dropoff"
	AtDropOff              = "at_dropoff"
	Delivered              = "delivered"
	Cancelled              = "cancelled"

	// electron status types
	AwaitingDispatch = "awaiting_dispatch"
	Dispatched       = "dispatched"
	OnTrip           = "on_trip"
	Offline          = "offline"

	// vehicle status
	Inactive = "inactive"

	// customer status
	Searching = "searching"
	AtRest    = "atRest"
)

type Delivery struct {
	Model
	Status               Status    `json:"status"`
	OriginLongitude      float64   `json:"originLongitude"`
	OriginLatitude       float64   `json:"originLatitude"`
	DestinationLongitude float64   `json:"destinationLongitude"`
	DestinationLatitude  float64   `json:"destinationLatitude"`
	Rating               float64   `json:"rating"`
	FinalCost            float64   `json:"finalCost"`
	InitialCost          float64   `json:"initialCost"`
	Completed            bool      `json:"completed"`
	Notes                string    `json:"notes"`
	CustomerID           uint      `json:"customerId"`
	Customer             Customer  `json:"customer"`
	ProductID            uint      `json:"productId"`
	Product              Product   `json:"product"`
	ElectronID           *uint     `json:"electronId,omitempty"`
	Electron             *Electron `json:"electron,omitempty"`
}
