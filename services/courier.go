package services

import (
	"github.com/gwuah/api/shared"
)

type DistanceFromOrigin float64
type DurationFromOrigin float64

func (s *Services) getDistanceAndDuration(origin shared.Coord, destination shared.Coord) (DurationFromOrigin, DistanceFromOrigin, error) {

	response, err := s.eta.GetDistanceFromOriginsToDestination([]shared.Coord{origin}, destination)

	if response.Code != "Ok" {
		return 0, 0, shared.MAPBOX_ERROR
	}

	if err != nil {
		return 0, 0, err
	}

	return DurationFromOrigin(*response.Durations[1][0]), DistanceFromOrigin(*response.Distances[1][0]), nil

}
