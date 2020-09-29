package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/lib/ws"
	"github.com/gwuah/api/shared"
	"github.com/ryankurte/go-mapbox/lib/base"
)

var (
	MAPBOX_ERROR = errors.New("Mapbox Request Failed")
)

func (s *Services) DispatchOrder(data shared.DeliveryRequest, order *models.Order, ws *ws.WSConnection) error {

	delivery := order.Deliveries[0]
	ids := s.GetClosestElectrons(shared.Coord{
		Latitude:  delivery.OriginLatitude,
		Longitude: delivery.OriginLongitude,
	}, 2)

	if len(ids) == 0 {
		log.Println("There are no drivers available")
		ws.SendMessage([]byte("There are no drivers available."))
		return nil
	}

	electrons, err := s.GetAllElectrons(ids)

	if err != nil {
		return err
	}

	coords := []base.Location{}

	for _, electron := range electrons {
		coords = append(coords, base.Location{
			Latitude:  electron.Latitude,
			Longitude: electron.Longitude,
		})
	}

	response := s.eta.GetDurationFromOrigin(base.Location{
		Latitude:  data.Origin.Latitude,
		Longitude: data.Origin.Longitude,
	}, coords)

	if response.Code != "Ok" {
		return MAPBOX_ERROR
	}

	durations := response.Durations
	var e []shared.ElectronWithEta

	for key, durationFromOrigin := range durations[0][1:] {
		electron := electrons[key]
		e = append(e, shared.ElectronWithEta{
			Electron: electron,
			Duration: durationFromOrigin,
		})
	}

	sort.Slice(e, func(i, j int) bool {
		return e[i].Duration < e[j].Duration
	})

	d := shared.NewOrder{
		Meta: shared.Meta{
			Type: "NewOrder",
		},
		Order: order,
	}

	convertedValue, err := json.Marshal(d)

	if err != nil {
		return nil
	}

	ticker := time.NewTicker(5 * time.Second)

electronLoop:
	for _, electron := range e {
		conn := s.hub.GetClient(fmt.Sprintf("electron_%s", electron.Electron.Id))
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

	fmt.Println("HANLDED")
	return nil

}
