package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/go-redis/redis"
	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/lib/ws"
	"github.com/gwuah/api/shared"
	"github.com/gwuah/api/utils/geo"
	"github.com/ryankurte/go-mapbox/lib/base"
)

var (
	MAPBOX_ERROR = errors.New("Mapbox Request Failed")
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
			log.Printf("Failed to load electrons in electron_index %d", index)
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
			log.Println("Failed To Load User", err)
		}
		electrons = append(electrons, electron)
	}

	return electrons, nil
}

func (s *Services) DispatchDelivery(data shared.DeliveryRequest, delivery *models.Delivery, ws *ws.WSConnection) error {

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

	delivery, err = s.repo.FindDelivery(delivery.ID)
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
