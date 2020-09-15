package dispatch

// This is the core of athena.
// Matching and Trip Management will be handled here

import (
	"fmt"

	"github.com/gwuah/api/wss"
)

type BaseMessage struct {
	Type string `json:"type"`
}

type Coord struct {
	Longitude float64 `json:"longitude" validate:"required"`
	Latitude  float64 `json:"latitude" validate:"required"`
}

type DeliveryRequest struct {
	BaseMessage BaseMessage `json:"baseMessage"`
	Origin      Coord       `json:"origin"`
	Destination Coord       `json:"destination"`
	ProductId   uint        `json:"productId"`
	Notes       string      `json:"notes"`
	CustomerID  uint        `json:"customerId"`
}

var MESSAGE_TYPES = map[string]string{
	"DeliveryRequest": "DeliveryRequest",
	"CancelDelivery":  "CancelDelivery",
	"GetEstimate":     "GetEstimate",
}

type Dispatch struct {
	maxMessageTypeLength int
}

func New() *Dispatch {
	return &Dispatch{
		maxMessageTypeLength: getLength(MESSAGE_TYPES["DeliveryRequest"]),
	}
}

func getLength(s string) int {
	return len(s)
}

func (d *Dispatch) getTypeOfMessage(message []byte) []byte {
	// this method pre-parses the message and extracts the type of message from the payload
	// this is done to speedup parsing and reduce size of marshalled/unmarshalled payload
	// custom algorithm, ask @gwuah for explanation
	start := 16
	end := start + d.maxMessageTypeLength + 2
	head := message[start:end]

	for {
		length := len(head)
		if length > 0 && head[length-1] != byte('"') {
			head = head[:length-1]
		} else {
			break
		}
	}

	return head

}

func (d *Dispatch) processIncomingMessage(message []byte, ws *wss.WSConnection) {
	// we parse the incoming message and parse it to the right handler to handle it
	switch string(d.getTypeOfMessage(message)) {
	case MESSAGE_TYPES["DeliveryRequest"]:
		fmt.Print("New Delivery Request")
	case MESSAGE_TYPES["CancelDelivery"]:
		fmt.Print("Cancel Delivery Request")
	case MESSAGE_TYPES["GetEstimate"]:
		fmt.Print("Get Estimate Request")
	}
}
