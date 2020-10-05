package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gwuah/api/lib/ws"
	"github.com/gwuah/api/shared"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var MESSAGE_TYPES = map[string]string{
	"DeliveryRequest":     "DeliveryRequest",
	"CancelDelivery":      "CancelDelivery",
	"GetEstimate":         "GetEstimate",
	"LocationUpdate":      "LocationUpdate",
	"GetClosestElectrons": "GetClosestElectrons",
	"AcceptDelivery":      "AcceptDelivery",
}

func (h *Handler) handleConnection(entity string) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		// this is unsafe, in future we have to set a static list of accepted origins
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

		if err != nil {
			log.Println("failed to setup websocket conn ..", err)
			return
		}

		wsConnection := &ws.WSConnection{
			Hub:                   h.Hub,
			Send:                  make(chan []byte),
			Conn:                  conn,
			Id:                    id,
			Entity:                entity,
			ProcessMessage:        h.processIncomingMessage,
			IsActive:              true,
			DeliveryAcceptanceAck: make(chan bool),
		}

		h.Hub.Register <- wsConnection

		go wsConnection.ReadPump()
		go wsConnection.WritePump()

	}
}

func (h *Handler) processIncomingMessage(message []byte, ws *ws.WSConnection) {
	switch string(h.getTypeOfMessage(message)) {
	case MESSAGE_TYPES["DeliveryRequest"]:
		h.processDeliveryRequest(message, ws)
	case MESSAGE_TYPES["CancelDelivery"]:
		h.handleDeliveryCancellation(message, ws)
	case MESSAGE_TYPES["LocationUpdate"]:
		h.handleLocationUpdate(message, ws)
	case MESSAGE_TYPES["GetClosestElectrons"]:
		h.handleGetClosestElectrons(message, ws)
	case MESSAGE_TYPES["AcceptDelivery"]:
		h.acceptDelivery(message, ws)
	default:
		log.Printf("No handler available for request %s", h.getTypeOfMessage(message))
	}
}

func (h *Handler) getTypeOfMessage(message []byte) []byte {
	// this method pre-parses the message and extracts the type of message from the payload
	// this is done to speedup parsing and reduce size of marshalled/unmarshalled payload
	// custom algorithm, ask @gwuah for explanation
	start := 16
	end := start + h.maxMessageTypeLength + 2
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

func (h *Handler) handleLocationUpdate(message []byte, ws *ws.WSConnection) {
	var data shared.UserLocationUpdate
	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Println("failed to parse message", err)
		return
	}

	err = h.Services.HandleLocationUpdate(data)
	if err != nil {
		log.Println("failed to handle location update", err)
		return
	}
}
