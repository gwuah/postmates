package handler

import (
	"encoding/json"
	"log"

	"github.com/gwuah/api/lib/ws"
	"github.com/gwuah/api/shared"
)

type ElectronWithEta struct {
	Electron *shared.User
	Duration float64
}

func (h *Handler) handleAcceptOrder(message []byte, ws *ws.WSConnection) {

	var data shared.AcceptOrder
	err := json.Unmarshal(message, &data)

	if err != nil {
		log.Println("Failed to parse message", err)
		return
	}

	err = h.Services.AcceptOrder(data, ws)
	if err != nil {
		log.Println("Failed to accept order", err)
		return
	}

}

func (h *Handler) handleDeliveryRequest(message []byte, ws *ws.WSConnection) {
	var data shared.DeliveryRequest

	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Println("Failed to parse message", err)
		return
	}

	product, err := h.Repo.FindProduct(data.ProductId)
	if err != nil {
		log.Printf("Failed to find product with id (%d)", data.ProductId)
		log.Println(err)
		return
	}

	if product.Name == "express" {

		_, order, err := h.Services.CreateDelivery(data)
		if err != nil {
			log.Println("failed to create delivery", err)
			return
		}

		err = h.Services.DispatchOrder(data, order, ws)
		if err != nil {
			log.Println("failed to dispatch order", err)
			return
		}

	} else {

	}
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
