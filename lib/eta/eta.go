package eta

import (
	"log"
	"os"

	mapbox "github.com/ryankurte/go-mapbox/lib"
	"github.com/ryankurte/go-mapbox/lib/base"
	directionsmatrix "github.com/ryankurte/go-mapbox/lib/directions_matrix"
)

func GetETAFromOrigin(origin base.Location, destinations []base.Location) *directionsmatrix.DirectionMatrixResponse {
	var token = os.Getenv("MAPBOX_TOKEN")
	mapBox, err := mapbox.NewMapbox(token)
	if err != nil {
		log.Println("Failed to initialize mapbox", err)
		return nil
	}

	var opts directionsmatrix.RequestOpts
	opts.SetSources([]string{"0"})
	opts.SetDestinations([]string{"all"})

	destinations = append([]base.Location{origin}, destinations...)

	response, err := mapBox.DirectionsMatrix.GetDirectionsMatrix(destinations, directionsmatrix.RoutingDriving, &opts)
	if err != nil {
		log.Println("Failed to get ETA", err)
	}

	return response

}
