package services

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/go-redis/redis"
	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/lib/ws"
	"github.com/gwuah/api/shared"
	"github.com/gwuah/api/utils/geo"
)

func (s *Services) IndexElectronLocation(param shared.UserLocationUpdate) (*shared.User, error) {
	newIndex := geo.CoordToIndex(param.Coord)

	electron, err := s.repo.GetElectronFromRedis(param.Id)

	if err == redis.Nil {
		electron = &shared.User{
			Id: param.Id,
		}
	}

	if err != redis.Nil && err != nil {
		return nil, err
	}

	oldIndex := electron.LastKnownIndex

	electron.Coord = param.Coord
	electron.LastKnownIndex = newIndex

	err = s.repo.InsertElectronIntoRedis(electron)

	if err != nil {
		return nil, err
	}

	if oldIndex != newIndex {
		err = s.repo.RemoveElectronFromIndex(oldIndex, electron)
		if err != nil {
			return nil, err
		}

		err = s.repo.InsertElectronIntoIndex(newIndex, electron)
		if err != nil {
			return nil, err
		}

	}

	return electron, nil
}

func (s *Services) GetClosestElectrons(coord shared.Coord, steps int) []string {

	rings := geo.GetRingsFromOrigin(coord, steps)

	electronsIds := []string{}

	for _, index := range rings {
		ids, err := s.repo.GetElectronsInIndex(index)

		if err != nil {
			log.Printf("failed to load electrons in electron_index %d", index)
			continue
		}

		if len(ids) > 0 {
			electronsIds = append(electronsIds, ids...)
		}
	}

	return electronsIds

}

func (s *Services) GetAllElectrons(ids []string) ([]*shared.User, error) {
	electrons := []*shared.User{}

	for _, id := range ids {
		electron, err := s.repo.GetElectronFromRedis(id)
		if err != nil {
			log.Println("failed To Load User", err)
		}
		electrons = append(electrons, electron)
	}

	return electrons, nil
}

func (s *Services) DispatchDelivery(data shared.DeliveryRequest, delivery *models.Delivery, ws *ws.WSConnection) error {

	foundElectronForOrder := false

dispatchLogic:

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

	coords := []shared.Coord{}

	for _, electron := range electrons {
		coords = append(coords, shared.Coord{
			Latitude:  electron.Latitude,
			Longitude: electron.Longitude,
		})
	}

	response, err := s.eta.GetDistanceFromOriginsToDestination(coords, shared.Coord{
		Latitude:  data.Origin.Latitude,
		Longitude: data.Origin.Longitude,
	})

	if err != nil {
		log.Println("mapbox request failed")
		return err
	}

	if response.Code != "Ok" {
		return shared.MAPBOX_ERROR
	}

	durations := response.Durations
	var e []shared.ElectronWithEta

	for key, duration := range durations[1:] {
		electron := electrons[key]
		e = append(e, shared.ElectronWithEta{
			Electron: electron,
			Duration: *duration[0],
		})
	}

	sort.Slice(e, func(i, j int) bool {
		return e[i].Duration < e[j].Duration
	})

	delivery, err = s.repo.FindDelivery(delivery.ID, true)
	if err != nil {
		return nil
	}

	d := shared.NewDelivery{
		Meta: shared.Meta{
			Type: "NewDelivery",
		},
		Delivery: delivery,
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
			continue
		}

		conn.SendMessage(convertedValue)

		select {
		case <-ticker.C:
			// move to next electron in queue
		case <-conn.DeliveryAcceptanceAck:
			// delivery has been accepted, exit
			ticker.Stop()
			foundElectronForOrder = true
			break electronLoop
		}

	}

	if !foundElectronForOrder {
		goto dispatchLogic
	}

	return nil
}
