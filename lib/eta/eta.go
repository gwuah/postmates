package eta

import (
	"log"
	"os"

	mapbox "github.com/ryankurte/go-mapbox/lib"
	"github.com/ryankurte/go-mapbox/lib/base"
	directionsmatrix "github.com/ryankurte/go-mapbox/lib/directions_matrix"
)

type Eta struct {
	token  string
	mapBox *mapbox.Mapbox
}

func New(token string) *Eta {
	if token == "" {
		log.Fatal("mapbox token required")
	}
	mapBox, err := mapbox.NewMapbox(token)

	if err != nil {
		log.Fatal("failed to initialize mapbox", err)
	}

	return &Eta{
		token:  token,
		mapBox: mapBox,
	}
}

func (e *Eta) _getDurationFromOrigin(origin base.Location, destinations []base.Location) *directionsmatrix.DirectionMatrixResponse {

	var opts directionsmatrix.RequestOpts
	opts.SetSources([]string{"0"})
	opts.SetDestinations([]string{"all"})

	destinations = append([]base.Location{origin}, destinations...)

	response, err := e.mapBox.DirectionsMatrix.GetDirectionsMatrix(destinations, directionsmatrix.RoutingDriving, &opts)
	if err != nil {
		log.Println("failed to get ETA", err)
	}

	return response

}

func (e *Eta) test_getDurationFromOrigin(origin base.Location, destinations []base.Location) *directionsmatrix.DirectionMatrixResponse {

	return &directionsmatrix.DirectionMatrixResponse{
		Code:      "Ok",
		Durations: [][]float64{{0, 10, 20, 21, 44}},
	}

}

func (e *Eta) GetDurationFromOrigin(origin base.Location, destinations []base.Location) *directionsmatrix.DirectionMatrixResponse {
	if os.Getenv("MODE") == "testing" {
		return e.test_getDurationFromOrigin(origin, destinations)
	}

	return e._getDurationFromOrigin(origin, destinations)

}
