package repository

import (
	"github.com/gwuah/api/database/models"
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
