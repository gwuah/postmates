package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/gwuah/api/lib/ws"
	"github.com/gwuah/api/shared"
	"github.com/ryankurte/go-mapbox/lib/base"
)

type ElectronWithEta struct {
	Electron *shared.User
	Duration float64
}

func (h *Handler) handleAcceptDeliveryRequest(message []byte, ws *ws.WSConnection) {

	var data shared.AcceptDeliveryRequest
	err := json.Unmarshal(message, &data)

	if err != nil {
		log.Println("Failed to parse message", err)
		return
	}

	delivery, err := h.Repo.FindDelivery(data.DeliveryId)

	if err != nil {
		log.Println("Failed to retrieve product with id", data.DeliveryId)
		return
	}

	electron := h.Hub.GetClient(fmt.Sprintf("electron_%s", ws.Id))
	if electron == nil {
		log.Println("Electron has disconnected from server", ws.Id)
		return
	}

	go func() {
		electron.AcceptDeliveryRequest([]byte("Trip Accepted"))
	}()

	customer := h.Hub.GetClient(fmt.Sprintf("customer_%d", delivery.CustomerID))
	if customer == nil {
		log.Println("Customer has disconnected from server")
		return
	}

	go func() {
		customer.SendMessage([]byte(fmt.Sprintf("Trip Accepted by electron %s", electron.Id)))
	}()

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

		response := h.Eta.GetDurationFromOrigin(base.Location{
			Latitude:  data.Origin.Latitude,
			Longitude: data.Origin.Longitude,
		}, coords)

		if response.Code != "Ok" {
			log.Println("Mapbox Request Failed", response.Code)
			return
		}

		durations := response.Durations
		var e []ElectronWithEta

		for key, durationFromOrigin := range durations[0][1:] {
			electron := electrons[key]
			e = append(e, ElectronWithEta{
				Electron: electron,
				Duration: durationFromOrigin,
			})
		}

		sort.Slice(e, func(i, j int) bool {
			return e[i].Duration < e[j].Duration
		})

		data := shared.NewDeliveryOrder{
			Meta: shared.Meta{
				Type: "NewDeliveryOrder",
			},
			Delivery: delivery,
		}

		convertedValue, err := json.Marshal(data)

		if err != nil {
			log.Println("Failed to marshal message", err)
			return
		}

		ticker := time.NewTicker(5 * time.Second)

	electronLoop:
		for _, electron := range e {
			conn := h.Hub.GetClient(fmt.Sprintf("electron_%s", electron.Electron.Id))
			if conn == nil {
				log.Printf("Electron %s has disconnected from server", electron.Electron.Id)
				continue
			}

			conn.SendMessage(convertedValue)

			select {
			case <-ticker.C:
				// move to next electron in queue
			case <-conn.AcceptDeliveryPipeline:
				// delivery has been accepted, exit
				ticker.Stop()
				break electronLoop
			}

		}

		fmt.Println("HANLDED == 100%")

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
