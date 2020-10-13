package repository

import (
	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/shared"
)

func (r *Repository) CreateTripPoint(data shared.UserLocationUpdate) (*models.TripPoint, error) {
	tripPoint := models.TripPoint{
		Latitude:   data.Latitude,
		Longitude:  data.Longitude,
		DeliveryID: data.DeliveryId,
		State:      data.State,
	}

	if err := r.DB.Create(&tripPoint).Error; err != nil {
		return nil, err
	}

	return &tripPoint, nil
}
