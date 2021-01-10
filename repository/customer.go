package repository

import (
	"github.com/gwuah/postmates/database/models"
)

func (r *Repository) FindCustomerByQuery(query string) (*models.Customer, error) {

	var customer models.Customer

	if err := r.DB.Where(query).First(&customer).Error; err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *Repository) FindCustomerByPhone(phone string) (*models.Customer, error) {

	customer := models.Customer{}
	err := r.DB.Where("phone = ?", phone).First(&customer).Error

	if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *Repository) CreateCustomerWithPhoneAndCode(phone string, code int) (*models.Customer, error) {

	customer := models.Customer{Phone: phone, Code: code}

	if err := r.DB.Create(&customer).Error; err != nil {
		return nil, err
	}

	return &customer, nil
}

func (r *Repository) UpdateCustomer(id uint, data map[string]interface{}) (*models.Customer, error) {
	var customer models.Customer

	if err := r.DB.Model(&customer).Where("id = ?", id).Updates(data).Error; err != nil {
		return nil, err
	}

	return &customer, nil
}
