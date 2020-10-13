package eta

import (
	"context"
	"log"
	"time"

	mapbox "github.com/airspacetechnologies/go-mapbox"
	"github.com/gwuah/api/shared"
)

type Eta struct {
	token  string
	mapBox *mapbox.Client
}

type DistanceFromOrigin float64
type DurationFromOrigin float64

func New(APIKey string) *Eta {
	if APIKey == "" {
		log.Fatal("mapbox token required")
	}

	mapboxClient, err := mapbox.NewClient(&mapbox.MapboxConfig{
		Timeout: 30 * time.Second,
		APIKey:  APIKey,
	})

	if err != nil {
		log.Fatal("failed to initialize mapbox", err)
	}

	return &Eta{
		token:  APIKey,
		mapBox: mapboxClient,
	}
}

func transformCoordinates(coords []shared.Coord) mapbox.Coordinates {
	transformed := mapbox.Coordinates{}

	for _, value := range coords {
		transformed = append(transformed, mapbox.Coordinate{
			Lat: value.Latitude,
			Lng: value.Longitude,
		})
	}

	return transformed
}

func (eta *Eta) GetDistanceFromOriginsToDestination(origins []shared.Coord, destination shared.Coord) (*mapbox.DirectionsMatrixResponse, error) {
	coords := mapbox.Coordinates{
		mapbox.Coordinate{
			Lat: destination.Latitude,
			Lng: destination.Longitude,
		},
	}
	coords = append(coords, transformCoordinates(origins)...)

	request := &mapbox.DirectionsMatrixRequest{
		Profile:       mapbox.ProfileDrivingTraffic,
		Coordinates:   coords,
		Annotations:   mapbox.Annotations{mapbox.AnnotationDistance, mapbox.AnnotationDuration},
		Destinations:  mapbox.Destinations{0},
		FallbackSpeed: 60,
	}

	response, err := eta.mapBox.DirectionsMatrix(context.TODO(), request)

	return response, err

}

func (eta *Eta) GetDistanceAndDuration(origin shared.Coord, destination shared.Coord) (DurationFromOrigin, DistanceFromOrigin, error) {

	response, err := eta.GetDistanceFromOriginsToDestination([]shared.Coord{origin}, destination)

	if response.Code != "Ok" {
		return 0, 0, shared.MAPBOX_ERROR
	}

	if err != nil {
		return 0, 0, err
	}

	return DurationFromOrigin(*response.Durations[1][0]), DistanceFromOrigin(*response.Distances[1][0]), nil

}
