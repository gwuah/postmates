package services

import (
	"github.com/gwuah/api/repository"
)

type Services struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Services {
	return &Services{repo: repo}
}
