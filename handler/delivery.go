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

	product, err := h.Repo.FindProduct(data.ProductId)
	if err != nil {
		log.Printf("Failed to find product with id (%d)", data.ProductId)
		log.Println(err)
		return
	}

	if product.Name == "express" {
		order, err := h.Repo.CreateOrder()
		if err != nil {
			log.Println("Failed to create order", err)
			return
		}

		rawDelivery, err := h.Repo.CreateDelivery(data, order)

		if err != nil {
			log.Println("Failed to create delivery", err)
			return
		}

		delivery, err := h.Repo.GetDelivery(rawDelivery.ID)

		if err != nil {
			log.Println("Failed to retrieve delivery", err)
			return
		}

		stringifiedResponse, err := json.Marshal(delivery)

		if err != nil {
			log.Println("Failed to marshal message", err)
			return
		}

		ws.SendMessage([]byte(stringifiedResponse))

	} else {

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
