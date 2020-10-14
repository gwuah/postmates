package services

import (
	"github.com/electra-systems/core-api/lib/billing"
	"github.com/electra-systems/core-api/lib/eta"
	"github.com/electra-systems/core-api/lib/ws"
	"github.com/electra-systems/core-api/repository"
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
