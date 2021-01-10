package eta

import (
	"context"
	"fmt"
	"log"

	"github.com/gwuah/postmates/shared"

	"googlemaps.github.io/maps"
)

type Eta struct {
	token string
	gmaps *maps.Client
}

type DistanceFromOrigin float64
type DurationFromOrigin float64

func New(googleAPIKey string) *Eta {

	if googleAPIKey == "" {
		log.Fatal("gmaps token required")
	}

	gmapsClient, err := maps.NewClient(maps.WithAPIKey(googleAPIKey))
	if err != nil {
		log.Fatal("failed to initialize gmaps", err)
	}

	return &Eta{
		token: googleAPIKey,
		gmaps: gmapsClient,
	}
}

func (eta *Eta) GMAPS__distanceMatrixBase(origins []shared.Coord, destinations []shared.Coord) (*maps.DistanceMatrixResponse, error) {

	modifiedOrigins := []string{}
	modifiedDestinations := []string{}

	for _, coord := range origins {
		modifiedOrigins = append(modifiedOrigins, fmt.Sprintf("%f,%f", coord.Latitude, coord.Longitude))
	}

	for _, coord := range destinations {
		modifiedDestinations = append(modifiedDestinations, fmt.Sprintf("%f,%f", coord.Latitude, coord.Longitude))
	}

	r := &maps.DistanceMatrixRequest{
		Origins:      modifiedOrigins,
		Destinations: modifiedDestinations,
	}

	resp, err := eta.gmaps.DistanceMatrix(context.Background(), r)

	if err != nil {
		return nil, err
	}

	return resp, nil

}

func (eta *Eta) GMAPS__getDistanceAndDuration1to1(origin shared.Coord, destination shared.Coord) (int, float64, error) {

	resp, err := eta.GMAPS__distanceMatrixBase([]shared.Coord{origin}, []shared.Coord{destination})

	if err != nil {
		return 0, 0, err
	}

	distance := resp.Rows[0].Elements[0].Distance.Meters
	duration := resp.Rows[0].Elements[0].Duration.Minutes()

	return distance, duration, nil

}

func (eta *Eta) GMAPS__getDistanceAndDurationManyTo1(origins []shared.Coord, destination shared.Coord) (*maps.DistanceMatrixResponse, error) {

	resp, err := eta.GMAPS__distanceMatrixBase(origins, []shared.Coord{destination})

	if err != nil {
		return nil, err
	}

	return resp, nil
}
