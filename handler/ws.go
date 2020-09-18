package handler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gwuah/api/lib/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var MESSAGE_TYPES = map[string]string{
	"DeliveryRequest":       "DeliveryRequest",
	"CancelDelivery":        "CancelDelivery",
	"GetEstimate":           "GetEstimate",
	"IndexElectronLocation": "IndexElectronLocation",
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

func (h *Handler) processIncomingMessage(message []byte, ws *ws.WSConnection) {
	switch string(h.getTypeOfMessage(message)) {
	case MESSAGE_TYPES["DeliveryRequest"]:
		h.handleDeliveryRequest(message, ws)
	case MESSAGE_TYPES["CancelDelivery"]:
		h.handleDeliveryCancellation(message, ws)
	case MESSAGE_TYPES["IndexElectronLocation"]:
		h.handleElectronLocationUpdate(message, ws)
	default:
		log.Printf("No handler available for request %s", h.getTypeOfMessage(message))
	}
}

func (h *Handler) handleConnection(entity string) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

		if err != nil {
			log.Println("Failed to setup websocket conn ..", err)
			return
		}

		wsConnection := &ws.WSConnection{
			Hub:            h.Hub,
			Send:           make(chan []byte),
			Conn:           conn,
			Id:             id,
			Entity:         entity,
			ProcessMessage: h.processIncomingMessage,
		}

		h.Hub.Register <- wsConnection

		go wsConnection.ReadPump()
		go wsConnection.WritePump()

	}
}
