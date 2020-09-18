package repository

import (
	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/shared"
	"gorm.io/gorm"
)

func (r *Repository) GetDelivery(id uint64) (models.Delivery, *gorm.DB) {

	var delivery models.Delivery

	result := r.DB.Find(&delivery, id)

	r.DB.Model(&delivery).Association("Customer").Find(&delivery.Customer)
	r.DB.Model(&delivery).Association("Order").Find(&delivery.Order)
	r.DB.Model(&delivery).Association("Product").Find(&delivery.Product)

	return delivery, result
}

func (r *Repository) CreateDelivery(data shared.DeliveryRequest, order *models.Order) (*models.Delivery, error) {
	delivery := models.Delivery{
		OriginLatitude:       data.Origin.Lat,
		OriginLongitude:      data.Origin.Lng,
		DestinationLatitude:  data.Destination.Lat,
		DestinationLongitude: data.Destination.Lng,
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
