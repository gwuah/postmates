package handler

import (
	"encoding/json"
	"log"

	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/lib/ws"
	"github.com/gwuah/api/shared"
)

type ElectronWithEta struct {
	Electron *shared.User
	Duration float64
}

func (h *Handler) acceptDelivery(message []byte, ws *ws.WSConnection) {

	var data shared.AcceptDelivery
	err := json.Unmarshal(message, &data)

	if err != nil {
		log.Println("failed to parse message", err)
		return
	}

	err = h.Services.AcceptDelivery(data, ws)
	if err != nil {
		log.Println("failed to accept delivery", err)
		return
	}

}

func (h *Handler) processDeliveryRequest(message []byte, ws *ws.WSConnection) {
	var data shared.DeliveryRequest

	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Println("failed to parse message", err)
		return
	}

	_, err = h.Repo.UpdateCustomer(data.CustomerID, map[string]interface{}{
		"State": models.Searching,
	})

	if err != nil {
		log.Println("failed to update customer", err)
		return
	}

	product, err := h.Repo.FindProduct(data.ProductId)
	if err != nil {
		log.Printf("failed to find product with id (%d)", data.ProductId)
		log.Println(err)
		return
	}

	if product.Name == "express" {

		delivery, err := h.Services.CreateDelivery(data)
		if err != nil {
			log.Println("failed to create delivery", err)
			return
		}

		err = h.Services.DispatchDelivery(data, delivery, ws)
		if err != nil {
			log.Println("failed to dispatch delivery", err)
			return
		}

	} else {

	}
}

func (h *Handler) handleDeliveryCancellation(message []byte, ws *ws.WSConnection) {
	var data shared.CancelDeliveryRequest
	err := json.Unmarshal(message, &data)
	if err != nil {
		log.Println("failed to parse message", err)
		return

	}

	ws.SendMessage([]byte("Delivery Cancelled"))

}
