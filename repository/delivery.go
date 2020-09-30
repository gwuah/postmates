package repository

import (
	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/shared"
	"gorm.io/gorm/clause"
)

func (r *Repository) CreateDelivery(data shared.DeliveryRequest, order *models.Order) (*models.Delivery, error) {
	delivery := models.Delivery{
		OriginLatitude:       data.Origin.Latitude,
		OriginLongitude:      data.Origin.Longitude,
		DestinationLatitude:  data.Destination.Latitude,
		DestinationLongitude: data.Destination.Longitude,
		Notes:                data.Notes,
		OrderID:              order.ID,
		ProductID:            data.ProductId,
		CustomerID:           data.CustomerID,
		Status:               models.Pending,
	}

	if err := r.DB.Create(&delivery).Error; err != nil {
		return nil, err
	}

	return &delivery, nil
}

func (r *Repository) FindDelivery(id uint) (*models.Delivery, error) {

	var delivery models.Delivery

	if err := r.DB.First(&delivery, id).Error; err != nil {
		return nil, err
	}

	r.DB.Preload(clause.Associations).Find(&delivery)

	return &delivery, nil
}

func (r *Repository) UpdateDelivery(id uint, data map[string]interface{}) (*models.Delivery, error) {
	var delivery models.Delivery

	if err := r.DB.Model(&delivery).Where("id = ?", id).Updates(data).Error; err != nil {
		return nil, err
	}

	return &delivery, nil
}
