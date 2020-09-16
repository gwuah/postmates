package dispatch

// This is the handlers of athena.
// Matching and Trip Management will be handled here

import (
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gwuah/api/repository"
	"github.com/gwuah/api/services"
	"github.com/gwuah/api/shared"

	"gorm.io/gorm"
)

var MESSAGE_TYPES = map[string]string{
	"DeliveryRequest": "DeliveryRequest",
	"CancelDelivery":  "CancelDelivery",
	"GetEstimate":     "GetEstimate",
}

type Dispatch struct {
	maxMessageTypeLength int
	hub                  *Hub
	services             *services.Services
	repo                 *repository.Repository
}

func New(DB *gorm.DB) *Dispatch {
	repo := repository.New(DB)
	services := services.New(repo)

	hub := NewHub()
	go hub.run()

	return &Dispatch{
		repo:                 repo,
		hub:                  hub,
		services:             services,
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
		var data shared.DeliveryRequest
		err := json.Unmarshal(message, &data)
		if err != nil {
			log.Println("Failed to parse message")
			break
		}
		d.services.CreateNewDeliveryRequest(data)
	case MESSAGE_TYPES["CancelDelivery"]:
		var data shared.CancelDeliveryRequest
		err := json.Unmarshal(message, &data)
		if err != nil {
			log.Println("Failed to parse message", err)
		}
		d.services.CancelDelivery(data)
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
