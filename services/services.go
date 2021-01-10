package services

import (
	"github.com/gwuah/postmates/lib/billing"
	"github.com/gwuah/postmates/lib/eta"
	"github.com/gwuah/postmates/lib/ws"
	"github.com/gwuah/postmates/repository"
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
