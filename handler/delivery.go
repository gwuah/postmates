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

		for _, electron := range e {
			fmt.Printf("%s is %dms away from request\n", electron.Electron.Id, int(electron.Duration))
			ws.SendMessage([]byte(fmt.Sprintf("Sending your request to %s, who is %dms away \n", electron.Electron.Id, int(electron.Duration))))

			conn := h.Hub.GetClient(fmt.Sprintf("electron_%s", electron.Electron.Id))
			if conn != nil {
				conn.SendMessage(convertedValue)
				fmt.Println("message sent")
			} else {
				fmt.Println("Electron has disconnected from server")

			}

			time.Sleep(2 * time.Second)

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
