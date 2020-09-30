package services

import (
	"fmt"

	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/lib/ws"
	"github.com/gwuah/api/shared"
)

func (s *Services) AcceptOrder(data shared.AcceptOrder, ws *ws.WSConnection) error {
	order, err := s.repo.UpdateOrder(data.OrderId, map[string]interface{}{
		"ElectronID": ws.Id,
	})
	if err != nil {
		return err
	}

	order, err = s.repo.FindOrder(data.OrderId)
	if err != nil {
		return err
	}

	_, err = s.repo.UpdateDelivery(order.Deliveries[0].ID, map[string]interface{}{
		"Status": models.PendingPickup,
	})
	if err != nil {
		return err
	}

	_, err = s.repo.UpdateElectron(order.ElectronID, map[string]interface{}{
		"Status": models.Dispatched,
	})
	if err != nil {
		return err
	}

	go func() {
		ws.AcceptDeliveryRequest([]byte("Trip Accepted"))
	}()

	for _, d := range order.Deliveries {
		customer := s.hub.GetClient(fmt.Sprintf("customer_%d", d.CustomerID))
		if customer != nil {
			go func() {
				customer.SendMessage([]byte(fmt.Sprintf("Trip Accepted by electron %s", ws.Id)))
			}()
		}
	}
	return nil

}
