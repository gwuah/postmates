// This is basically our data layer.
package repository

import "gorm.io/gorm"

type Repository struct {
	DB *gorm.DB
}

func New(db *gorm.DB) *Repository {
	return &Repository{db}
}
