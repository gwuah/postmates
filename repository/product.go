package repository

import (
	"errors"

	"github.com/electra-systems/core-api/database/models"
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

func (r *Repository) FindAllProducts() ([]models.Product, error) {

	products := []models.Product{}

	if err := r.DB.Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}
