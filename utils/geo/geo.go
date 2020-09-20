package geo

import (
	"github.com/gwuah/api/shared"
	"github.com/uber/h3-go"
)

func CoordToIndex(param shared.Coord) h3.H3Index {
	return h3.FromGeo(h3.GeoCoord{
		Latitude:  param.Lat,
		Longitude: param.Lng,
	}, 8)
}

func GetRingsFromOrigin(coord shared.Coord, steps int) []h3.H3Index {
	return h3.KRing(CoordToIndex(coord), steps)
}
