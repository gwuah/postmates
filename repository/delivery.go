package repository

import (
	"fmt"

	"github.com/gwuah/postmates/database/models"
	"github.com/gwuah/postmates/shared"
	"gorm.io/gorm/clause"
)

func (r *Repository) CreateDelivery(data shared.DeliveryRequest) (*models.Delivery, error) {
	delivery := models.Delivery{
		OriginLatitude:       data.Origin.Latitude,
		OriginLongitude:      data.Origin.Longitude,
		DestinationLatitude:  data.Destination.Latitude,
		DestinationLongitude: data.Destination.Longitude,
		Notes:                data.Notes,
		ProductID:            data.ProductId,
		CustomerID:           data.CustomerID,
		State:                models.Pending,
	}

	if err := r.DB.Create(&delivery).Error; err != nil {
		return nil, err
	}

	return &delivery, nil
}

func (r *Repository) FindDelivery(id uint, loadAssociations bool) (*models.Delivery, error) {

	var delivery models.Delivery

	if err := r.DB.First(&delivery, id).Error; err != nil {
		return nil, err
	}

	if loadAssociations {
		r.DB.Preload(clause.Associations).Find(&delivery)
	}

	return &delivery, nil
}

func (r *Repository) UpdateDelivery(id uint, data map[string]interface{}) (*models.Delivery, error) {
	var delivery models.Delivery

	if err := r.DB.Model(&delivery).Where("id = ?", id).Updates(data).Error; err != nil {
		return nil, err
	}

	return &delivery, nil
}

func (r *Repository) DeliveryCount(condition string) (int64, error) {
	var count int64
	if err := r.DB.Model(&models.Delivery{}).Where(condition).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository) DeliverySum(condition string, field string) (int64, error) {
	var count int64
	response := r.DB.Model(&models.Delivery{}).Where(condition).Select(fmt.Sprintf("sum(%s)", field))

	if err := response.Error; err != nil {
		return 0, err
	}

	response.Scan(&count)
	return count, nil
}
