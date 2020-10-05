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

func (s *Services) HandleLocationUpdate(params shared.UserLocationUpdate) error {

	switch params.State {
	case models.AwaitingDispatch:
		_, err := s.indexCourierLocation(params)
		return err
	case models.Dispatched, models.OnTrip:
		_, err := s.indexCourierLocation(params)
		if err != nil {
			return err
		}
		_, err = s.repo.CreateTripPoint(params)
		if err != nil {
			return err
		}
		err = s.relayCoordsToCustomer(params)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Services) relayCoordsToCustomer(params shared.UserLocationUpdate) error {
	delivery, err := s.repo.FindDelivery(params.DeliveryId, false)
	if err != nil {
		return err
	}

	customerConn := s.hub.GetClient(fmt.Sprintf("customer_%d", delivery.CustomerID))

	customerConn.SendMessage([]byte(""))

	return nil
}

func (s *Services) indexCourierLocation(param shared.UserLocationUpdate) (*shared.User, error) {
	newIndex := geo.CoordToIndex(param.Coord)

	courier, err := s.repo.GetCourierFromRedis(param.Id)

	if err == redis.Nil {
		courier = &shared.User{
			Id: param.Id,
		}
	}

	if err != redis.Nil && err != nil {
		return nil, err
	}

	oldIndex := courier.LastKnownIndex

	courier.Coord = param.Coord
	courier.LastKnownIndex = newIndex

	err = s.repo.InsertCourierIntoRedis(courier)

	if err != nil {
		return nil, err
	}

	if oldIndex != newIndex {
		err = s.repo.RemoveCourierFromIndex(oldIndex, courier)
		if err != nil {
			return nil, err
		}

		err = s.repo.InsertCourierIntoIndex(newIndex, courier)
		if err != nil {
			return nil, err
		}

	}

	return courier, nil
}

func (s *Services) GetClosestCouriers(coord shared.Coord, steps int) []string {

	rings := geo.GetRingsFromOrigin(coord, steps)

	couriersIds := []string{}

	for _, index := range rings {
		ids, err := s.repo.GetCouriersInIndex(index)

		if err != nil {
			log.Printf("failed to load couriers in courier_index %d", index)
			continue
		}

		if len(ids) > 0 {
			couriersIds = append(couriersIds, ids...)
		}
	}

	return couriersIds

}

func (s *Services) GetAllCouriers(ids []string) ([]*shared.User, error) {
	couriers := []*shared.User{}

	for _, id := range ids {
		courier, err := s.repo.GetCourierFromRedis(id)
		if err != nil {
			log.Println("failed To Load User", err)
		}
		couriers = append(couriers, courier)
	}

	return couriers, nil
}

func (s *Services) DispatchDelivery(data shared.DeliveryRequest, delivery *models.Delivery, ws *ws.WSConnection) error {

	foundCourierForOrder := false

dispatchLogic:

	ids := s.GetClosestCouriers(shared.Coord{
		Latitude:  delivery.OriginLatitude,
		Longitude: delivery.OriginLongitude,
	}, 2)

	if len(ids) == 0 {
		log.Println("There are no drivers available")
		ws.SendMessage([]byte("There are no drivers available."))
		return nil
	}

	couriers, err := s.GetAllCouriers(ids)

	if err != nil {
		return err
	}

	coords := []shared.Coord{}

	for _, courier := range couriers {
		coords = append(coords, shared.Coord{
			Latitude:  courier.Latitude,
			Longitude: courier.Longitude,
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
	var e []shared.CourierWithEta

	for key, duration := range durations[1:] {
		courier := couriers[key]
		e = append(e, shared.CourierWithEta{
			Courier:  courier,
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

	ticker := time.NewTicker(5 * time.Second)

courierLoop:
	for _, courier := range e {
		conn := s.hub.GetClient(fmt.Sprintf("courier_%s", courier.Courier.Id))
		if conn == nil {
			continue
		}

		duration, distance, err := s.getDistanceAndDuration(shared.Coord{
			Latitude:  courier.Courier.Latitude,
			Longitude: courier.Courier.Longitude,
		}, shared.Coord{
			Latitude:  delivery.OriginLatitude,
			Longitude: delivery.OriginLongitude,
		})

		convertedDeliveryRequest, err := json.Marshal(shared.NewDelivery{
			Meta: shared.Meta{
				Type: "NewDelivery",
			},
			Delivery:         delivery,
			DistanceToPickup: float64(distance),
			DurationToPickup: float64(duration),
		})

		if err != nil {
			return nil
		}

		conn.SendMessage(convertedDeliveryRequest)

		select {
		case <-ticker.C:
			// move to next courier in queue
		case <-conn.DeliveryAcceptanceAck:
			// delivery has been accepted, exit
			ticker.Stop()
			foundCourierForOrder = true
			break courierLoop
		}

	}

	if !foundCourierForOrder {
		goto dispatchLogic
	}

	return nil
}
