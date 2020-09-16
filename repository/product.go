package repository

import (
	"errors"

	"github.com/gwuah/api/database/models"
)

func (r *Repository) FindProduct(id uint) (*models.Product, error) {

	product := models.Product{}

	if err := r.DB.First(&product, id).Error; err != nil {
		return nil, err
	}

	if product.ID == 0 {
		return nil, errors.New("Product Doesn't Exist")
	}

	return &product, nil
}
