package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gwuah/api/lib/eta"
	"github.com/gwuah/api/lib/ws"
	"github.com/gwuah/api/shared"
	"github.com/ryankurte/go-mapbox/lib/base"
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

		ids := h.Services.GetClosestElectrons(shared.Coord{
			Latitude:  delivery.OriginLatitude,
			Longitude: delivery.OriginLongitude,
		}, 2)

		if len(ids) == 0 {
			log.Println("There are no drivers available")
			ws.SendMessage([]byte("There are no drivers available."))
			return
		}

		electrons, err := h.Services.GetAllElectrons(ids)

		if err != nil {
			log.Println("Failed to marshal message", err)
			return
		}

		coords := []base.Location{}

		for _, electron := range electrons {
			coords = append(coords, base.Location{
				Latitude:  electron.Latitude,
				Longitude: electron.Longitude,
			})
		}

		if err != nil {
			log.Println("Failed to mass get electrons", err)
			return
		}

		response := eta.GetETAFromOrigin(base.Location{
			Latitude:  data.Origin.Latitude,
			Longitude: data.Origin.Longitude,
		}, coords)

		if response.Code != "Ok" {
			log.Println("Mapbox Request Failed", response.Code)
			return
		}

		durations := response.Durations
		fmt.Println(durations)

		// mapIdToElectron

		for key, durationFromOrigin := range durations[0][1:] {
			electron := electrons[key]
			fmt.Printf("%s is %dms away from request\n", electron.Id, int(durationFromOrigin))
		}
		ws.SendMessage(stringifiedResponse)

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
