package repository

import (
	"errors"

	"github.com/gwuah/api/database/models"
	"gorm.io/gorm/clause"
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

func (r *Repository) FindOrder(id uint) (*models.Order, error) {

	order := models.Order{Electron: &models.Electron{}}

	if err := r.DB.First(&order, id).Error; err != nil {
		return nil, err
	}

	if order.ID == 0 {
		return nil, errors.New("order doesn't exist")
	}

	r.DB.Preload(clause.Associations).Find(&order)

	return &order, nil
}

func (r *Repository) UpdateOrder(id uint, data map[string]interface{}) (*models.Order, error) {
	order := models.Order{}

	if err := r.DB.Model(&order).Where("id = ?", id).Updates(data).Error; err != nil {
		return nil, err
	}

	return &order, nil
}
