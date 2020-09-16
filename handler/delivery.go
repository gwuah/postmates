package handler

import (
	"encoding/json"
	"log"

	"github.com/gwuah/api/lib/ws"
	"github.com/gwuah/api/shared"
)

func (h *Handler) handleDeliveryRequest(message []byte, ws *ws.WSConnection) {
	var data shared.DeliveryRequest
	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Println("Failed to parse message", err)
		return
	}

	ws.SendMessage([]byte("Delivery Placed"))

}

func (h *Handler) handleDeliveryCancellation(message []byte, ws *ws.WSConnection) {
	var data shared.CancelDeliveryRequest
	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Println("Failed to parse message", err)
		return

	}

	ws.SendMessage([]byte("Delivery Cancelled"))

}
