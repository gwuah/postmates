package repository

import (
	"github.com/gwuah/api/database/models"
	"gorm.io/gorm"
)

type CreateOrderSchema struct {
	Phone string `json:"phone" validate:"required"`
}

func (r *Repository) CreateOrder(data CreateOrderSchema) (models.Order, *gorm.DB) {

	order := models.Order{}

	result := r.DB.Create(&order)

	return order, result
}
