package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/electra-systems/core-api/database/models"
	"github.com/electra-systems/core-api/shared"
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

func (s *Services) RateDelivery(data shared.RatingRequest) (bool, error) {

	if data.IsCustomerRating {
		delivery, err := s.repo.UpdateDelivery(data.CustomerRating.DeliveryId, map[string]interface{}{
			"CustomerRating":        data.CustomerRating.Rating,
			"CustomerRatingMessage": data.CustomerRating.Message,
		})
		if err != nil {
			return false, err
		}

		delivery, err = s.repo.FindDelivery(data.CustomerRating.DeliveryId, false)
		if err != nil {
			return false, err
		}

		if delivery.State != models.Completed {
			return false, errors.New("delivery not completed")
		}

		condition := fmt.Sprintf("courier_id = %d AND state = %s", *delivery.CourierID, models.Completed)

		totalTrips, err := s.repo.DeliveryCount(condition)
		if err != nil {
			return false, err
		}

		totalRatings, err := s.repo.DeliverySum(condition, "customer_rating")
		if err != nil {
			return false, err
		}

		averageRating := totalRatings / totalTrips

		_, err = s.repo.UpdateCourier(*delivery.CourierID, map[string]interface{}{
			"Rating": averageRating,
		})

		if err != nil {
			return false, err
		}

	} else {

		delivery, err := s.repo.UpdateDelivery(data.CourierRating.DeliveryId, map[string]interface{}{
			"CourierRating":        data.CourierRating.Rating,
			"CourierRatingMessage": data.CourierRating.Message,
		})

		if err != nil {
			return false, err
		}

		delivery, err = s.repo.FindDelivery(data.CourierRating.DeliveryId, false)
		if err != nil {
			return false, err
		}

		condition := fmt.Sprintf("customer_id = %d", delivery.CustomerID)

		totalTrips, err := s.repo.DeliveryCount(condition)
		if err != nil {
			return false, err
		}

		totalRatings, err := s.repo.DeliverySum(condition, "courier_rating")
		if err != nil {
			return false, err
		}

		averageRating := totalRatings / totalTrips

		_, err = s.repo.UpdateCustomer(delivery.CustomerID, map[string]interface{}{
			"Rating": averageRating,
		})

		if err != nil {
			return false, err
		}

	}

	return true, nil
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
