package geo

import (
	"github.com/electra-systems/core-api/shared"
	"github.com/uber/h3-go"
)

func CoordToIndex(param shared.Coord) h3.H3Index {
	return h3.FromGeo(h3.GeoCoord{
		Latitude:  param.Latitude,
		Longitude: param.Longitude,
	}, 8)
}

func GetRingsFromOrigin(coord shared.Coord, steps int) []h3.H3Index {
	return h3.KRing(CoordToIndex(coord), steps)
}

func ConvertMetresToKM(distance float64) float64 {
	return distance / 1000
}
