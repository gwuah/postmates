package services

import (
	"encoding/json"
	"fmt"
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

func (s *Services) AcceptDelivery(data shared.AcceptDelivery, electronWS *ws.WSConnection) error {
	electronFromRedis, err := s.repo.GetElectronFromRedis(electronWS.Id)
	if err != nil {
		log.Printf("failed to retrieve electron %s from redis", electronWS.Id)
		return nil
	}

	delivery, err := s.repo.FindDelivery(data.DeliveryId, false)
	if err != nil {
		return err
	}

	duration, distance, err := s.getDistanceAndDuration(shared.Coord{
		Latitude:  electronFromRedis.Latitude,
		Longitude: electronFromRedis.Longitude,
	}, shared.Coord{
		Latitude:  delivery.OriginLatitude,
		Longitude: delivery.OriginLongitude,
	})
	if err != nil {
		log.Printf("failed to get electron %s ETA", electronWS.Id)
		return nil
	}

	_, err = s.repo.UpdateDelivery(data.DeliveryId, map[string]interface{}{
		"State":      models.PendingPickup,
		"ElectronID": electronWS.Id,
	})
	if err != nil {
		return err
	}

	_, err = s.repo.UpdateElectron(uint(utils.ConvertToUint64(electronWS.Id)), map[string]interface{}{
		"State": models.Dispatched,
	})
	if err != nil {
		return err
	}

	go func() {
		electronWS.AckDeliveryAcceptance(true)
	}()

	delivery, err = s.repo.FindDelivery(data.DeliveryId, true)
	if err != nil {
		return err
	}

	electron, err := s.repo.FindElectron(*delivery.ElectronID)
	if err != nil {
		return err
	}

	customer := s.hub.GetClient(fmt.Sprintf("customer_%d", delivery.CustomerID))
	if customer != nil {
		go func() {
			electron.Latitude = electronFromRedis.Latitude
			electron.Longitude = electronFromRedis.Longitude

			acceptanceDataStruct := shared.DeliveryAcceptedPayload{
				Meta: shared.Meta{
					Type: "DeliveryAccepted",
				},
				Electron:         *electron,
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
