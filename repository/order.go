package repository

import (
	"github.com/gwuah/api/database/models"
)

type CreateOrderSchema struct {
	Phone string `json:"phone" validate:"required"`
}

func (r *Repository) CreateOrder() (*models.Order, error) {
	order := models.Order{}

	if err := r.DB.Create(&order).Error; err != nil {
		return nil, err
	}

	return &order, nil

}
