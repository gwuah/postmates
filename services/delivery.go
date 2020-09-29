package services

import (
	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/shared"
)

func (s *Services) CreateDelivery(data shared.DeliveryRequest) (*models.Delivery, *models.Order, error) {
	order, err := s.repo.CreateOrder()

	if err != nil {
		return nil, nil, err
	}

	delivery, err := s.repo.CreateDelivery(data, order)

	if err != nil {
		return nil, nil, err
	}

	order, err = s.repo.FindOrder(order.ID)

	if err != nil {
		return nil, nil, err
	}

	return delivery, order, nil

}

func (s *Services) CancelDelivery(data shared.CancelDeliveryRequest) bool {
	return true
}
