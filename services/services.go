package services

import (
	"github.com/gwuah/api/lib/billing"
	"github.com/gwuah/api/lib/eta"
	"github.com/gwuah/api/lib/ws"
	"github.com/gwuah/api/repository"
)

type Services struct {
	repo    *repository.Repository
	eta     *eta.Eta
	hub     *ws.Hub
	billing *billing.Billing
}

func New(repo *repository.Repository, eta *eta.Eta, hub *ws.Hub, billing *billing.Billing) *Services {
	return &Services{repo: repo, eta: eta, hub: hub, billing: billing}
}
