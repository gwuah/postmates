package services

import (
	"strings"

	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/shared"
)

func (s *Services) CreateNewDeliveryRequest(data shared.DeliveryRequest) (*models.Delivery, error) {

	product, err := s.repo.FindProduct(data.ProductId)

	if err != nil {
		return nil, err
	}

	if strings.ToLower(product.Name) == "express" {
		order, err := s.repo.CreateOrder()

		if err != nil {
			return nil, err
		}

		delivery, err := s.repo.CreateDelivery(data, order)

		if err != nil {
			return nil, err
		}

		return delivery, nil

	} else if strings.ToLower(product.Name) == "pool" {

	}

	return nil, nil

}

func (s *Services) CancelDelivery(data shared.CancelDeliveryRequest) bool {
	return true
}
