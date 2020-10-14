package billing

import (
	"math"

	"github.com/electra-systems/core-api/utils/geo"
)

const (
	BASE_PRICE = 5
)

type Billing struct {
}

func New() *Billing {
	return &Billing{}
}

func (b *Billing) GetDeliveryCost(distance float64) float64 {
	fare := (13 * geo.ConvertMetresToKM(distance)) / 12.5
	if fare < BASE_PRICE {
		return BASE_PRICE
	}
	return math.Ceil(fare)
}
