package services

import (
	"encoding/json"
	"log"

	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/lib/ws"
	"github.com/gwuah/api/shared"
	"github.com/gwuah/api/utils"
)

func (s *Services) CreateDelivery(data shared.DeliveryRequest) (*models.Delivery, error) {
	delivery, err := s.repo.CreateDelivery(data)
	if err != nil {
		return nil, err
	}
	return delivery, nil
}

func (s *Services) AcceptDelivery(data shared.AcceptDelivery, courierWS *ws.WSConnection) error {
	courierFromRedis, err := s.repo.GetCourierFromRedis(courierWS.Id)
	if err != nil {
		log.Printf("failed to retrieve courier %s from redis", courierWS.Id)
		return nil
	}

	delivery, err := s.repo.FindDelivery(data.DeliveryId, false)
	if err != nil {
		return err
	}

	duration, distance, err := s.eta.GMAPS__getDistanceAndDuration1to1(shared.Coord{
		Latitude:  courierFromRedis.Latitude,
		Longitude: courierFromRedis.Longitude,
	}, shared.Coord{
		Latitude:  delivery.OriginLatitude,
		Longitude: delivery.OriginLongitude,
	})
	if err != nil {
		log.Printf("failed to get courier %s ETA", courierWS.Id)
		return nil
	}

	_, err = s.repo.UpdateDelivery(data.DeliveryId, map[string]interface{}{
		"State":     models.PendingPickup,
		"CourierID": courierWS.Id,
	})
	if err != nil {
		return err
	}

	_, err = s.repo.UpdateCourier(uint(utils.ConvertToUint64(courierWS.Id)), map[string]interface{}{
		"State": models.Dispatched,
	})
	if err != nil {
		return err
	}

	go func() {
		courierWS.AckDeliveryAcceptance(true)
	}()

	delivery, err = s.repo.FindDelivery(data.DeliveryId, true)
	if err != nil {
		return err
	}

	courier, err := s.repo.FindCourier(*delivery.CourierID)
	if err != nil {
		return err
	}

	customer := s.hub.GetCustomer(delivery.CustomerID)
	if customer != nil {
		go func() {
			courier.Latitude = courierFromRedis.Latitude
			courier.Longitude = courierFromRedis.Longitude

			acceptanceDataStruct := shared.DeliveryAcceptedPayload{
				Meta: shared.Meta{
					Type: "DeliveryAccepted",
				},
				Courier:          *courier,
				Delivery:         *delivery,
				DistanceToPickup: float64(distance),
				DurationToPickup: float64(duration),
			}

			acceptanceData, err := json.Marshal(acceptanceDataStruct)
			if err != nil {
				return
			}

			customer.SendMessage(acceptanceData)
		}()
	}
	return nil
}

func (s *Services) CancelDelivery(data shared.CancelDeliveryRequest) bool {
	return true
}
