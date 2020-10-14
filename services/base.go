package services

import (
	"errors"
	"strings"

	"github.com/gwuah/api/shared"
)

type PricePerProduct struct {
	ProductId uint    `json:"productId"`
	Price     float64 `json:"price"`
}

type GetDeliveryCostResponse struct {
	Estimates map[uint]PricePerProduct `json:"estimate"`
	Distance  float64                  `json:"distance"`
	Duration  float64                  `json:"duration"`
}

func (s *Services) GetDeliveryCost(data shared.GetDeliveryCostRequest) (*GetDeliveryCostResponse, error) {
	products, err := s.repo.FindAllProducts()

	if err != nil {

		return nil, errors.New("failed to load products")
	}

	duration, distance, err := s.eta.GMAPS__getDistanceAndDuration1to1(data.Origin, data.Destination)

	if err != nil {
		return nil, err
	}

	response := GetDeliveryCostResponse{
		Estimates: make(map[uint]PricePerProduct),
		Duration:  float64(duration),
		Distance:  float64(distance),
	}

	for _, product := range products {
		switch strings.ToLower(product.Name) {
		case "express":
			response.Estimates[product.ID] = PricePerProduct{
				ProductId: product.ID,
				Price:     s.billing.GetDeliveryCost(float64(distance)),
			}
			break

		}
	}

	return &response, nil

}
