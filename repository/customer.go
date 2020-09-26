package repository

import (
	"github.com/gwuah/api/database/models"
)

func (r *Repository) FindCustomerByPhone(phone string) (*models.Customer, error) {

	customer := models.Customer{}
	err := r.DB.Where("phone = ?", phone).First(&customer).Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *Repository) CreateCustomerWithPhone(phone string) (*models.Customer, error) {

	customer := models.Customer{Phone: phone}
	err := r.DB.Create(&customer).Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}
