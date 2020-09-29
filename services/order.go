package services

import (
	"errors"
	"fmt"

	"github.com/gwuah/api/lib/ws"
	"github.com/gwuah/api/shared"
)

func (s *Services) AcceptOrder(data shared.AcceptOrder, ws *ws.WSConnection) error {
	order, err := s.repo.FindAndUpdateOrder(data.OrderId, map[string]interface{}{
		"ElectronID": ws.Id,
	})

	if err != nil {
		return err
	}

	electron := s.hub.GetClient(ws.GetIdBasedOnType())
	if electron == nil {
		return errors.New("electron has disconnected from server")
	}

	go func() {
		electron.AcceptDeliveryRequest([]byte("Trip Accepted"))
	}()

	for _, d := range order.Deliveries {
		customer := s.hub.GetClient(fmt.Sprintf("customer_%d", d.CustomerID))
		if customer != nil {
			go func() {
				customer.SendMessage([]byte(fmt.Sprintf("Trip Accepted by electron %s", electron.Id)))
			}()
		}
	}
	return nil

}
