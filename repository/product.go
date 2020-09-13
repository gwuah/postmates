package repository

import (
	"github.com/gwuah/api/database/models"
	"gorm.io/gorm"
)

func (r *Repository) FindProduct(id uint) (models.Product, *gorm.DB) {

	product := models.Product{}

	result := r.DB.First(&product, id)

	return product, result
}
