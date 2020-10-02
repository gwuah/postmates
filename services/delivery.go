package services

import (
	"fmt"

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

func (s *Services) AcceptDelivery(data shared.AcceptDelivery, ws *ws.WSConnection) error {
	_, err := s.repo.UpdateDelivery(data.DeliveryId, map[string]interface{}{
		"Status":     models.PendingPickup,
		"ElectronID": ws.Id,
	})
	if err != nil {
		return err
	}

	_, err = s.repo.UpdateElectron(uint(utils.ConvertToUint64(ws.Id)), map[string]interface{}{
		"Status": models.Dispatched,
	})
	if err != nil {
		return err
	}

	go func() {
		ws.AcceptDeliveryRequest([]byte("Trip Accepted"))
	}()

	delivery, err := s.repo.FindDelivery(data.DeliveryId)
	customer := s.hub.GetClient(fmt.Sprintf("customer_%d", delivery.CustomerID))
	if customer != nil {
		go func() {
			customer.SendMessage([]byte(fmt.Sprintf("Trip Accepted by electron %s", ws.Id)))
		}()
	}
	return nil

}

func (s *Services) CancelDelivery(data shared.CancelDeliveryRequest) bool {
	return true
}
