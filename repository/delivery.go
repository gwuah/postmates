package repository

import (
	"errors"

	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/shared"
)

func (r *Repository) GetDelivery(id uint) (*models.Delivery, error) {

	var delivery models.Delivery

	err := r.DB.Find(&delivery, id).Error

	if err != nil {
		return nil, err
	}

	r.DB.Model(&delivery).Association("Customer").Find(&delivery.Customer)
	r.DB.Model(&delivery).Association("Order").Find(&delivery.Order)
	r.DB.Model(&delivery).Association("Product").Find(&delivery.Product)

	return &delivery, nil
}

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
	}

	if err := r.DB.Create(&delivery).Error; err != nil {
		return nil, err
	}

	return &delivery, nil
}

func (r *Repository) FindDelivery(id uint) (*models.Delivery, error) {

	delivery := models.Delivery{}

	if err := r.DB.First(&delivery, id).Error; err != nil {
		return nil, err
	}

	if delivery.ID == 0 {
		return nil, errors.New("Product Doesn't Exist")
	}

	return &delivery, nil
}
