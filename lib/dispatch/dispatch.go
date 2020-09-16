package dispatch

// This is the core of athena.
// Matching and Trip Management will be handled here

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

type BaseMessage struct {
	Type string `json:"type"`
}

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

type CancelDelivery struct {
	BaseMessage BaseMessage `json:"meta"`
	TripId      uint        `json:"tripId"`
}

var MESSAGE_TYPES = map[string]string{
	"DeliveryRequest": "DeliveryRequest",
	"CancelDelivery":  "CancelDelivery",
	"GetEstimate":     "GetEstimate",
}

type Dispatch struct {
	maxMessageTypeLength int
	hub                  *Hub
}

func New() *Dispatch {
	hub := NewHub()
	go hub.run()

	return &Dispatch{
		hub:                  hub,
		maxMessageTypeLength: len(MESSAGE_TYPES["DeliveryRequest"]),
	}
}

func (d *Dispatch) getTypeOfMessage(message []byte) []byte {
	// this method pre-parses the message and extracts the type of message from the payload
	// this is done to speedup parsing and reduce size of marshalled/unmarshalled payload
	// custom algorithm, ask @gwuah for explanation
	start := 16
	end := start + d.maxMessageTypeLength + 2
	payload := message[start:end]

	numberOfQuotesSeen := 0
	head := []byte{}

	for i := 0; i < len(payload); i++ {
		value := payload[i]

		if numberOfQuotesSeen == 2 {
			break
		}

		if value == byte('"') {
			numberOfQuotesSeen++
		}

		head = append(head, value)

	}

	return head[1 : len(head)-1]

}

func (d *Dispatch) processIncomingMessage(message []byte, ws *WSConnection) {
	switch string(d.getTypeOfMessage(message)) {
	case MESSAGE_TYPES["DeliveryRequest"]:
		ws.sendMessage([]byte("New Delivery Request"))
		var request DeliveryRequest
		err := json.Unmarshal(message, &request)
		if err != nil {
			log.Println("Failed to parse message")
		}
		log.Println(request)
	case MESSAGE_TYPES["CancelDelivery"]:
		ws.sendMessage([]byte("Cancel Delivery Request"))
		var request CancelDelivery
		err := json.Unmarshal(message, &request)
		if err != nil {
			log.Println("Failed to parse message", err)
		}
		log.Println(request)
	case MESSAGE_TYPES["GetEstimate"]:
		ws.sendMessage([]byte("Get Estimate Request"))
	default:
		log.Println("No Match")
	}
}

func (d *Dispatch) HandleConnection(entity string) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

		if err != nil {
			log.Println("Failed to setup websocket conn ..", err)
			return
		}

		wsConnection := &WSConnection{
			hub:            d.hub,
			send:           make(chan []byte),
			conn:           conn,
			id:             id,
			entity:         entity,
			processMessage: d.processIncomingMessage,
		}

		d.hub.register <- wsConnection

		go wsConnection.getIncomingMessages()
		go wsConnection.writeMessageToClient()

	}
}
