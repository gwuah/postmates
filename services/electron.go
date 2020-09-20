package services

import (
	"log"

	"github.com/go-redis/redis"
	"github.com/gwuah/api/shared"
	"github.com/gwuah/api/utils/geo"
)

func (s *Services) IndexElectronLocation(param shared.UserLocationUpdate) (*shared.User, error) {
	newIndex := geo.CoordToIndex(param.Coord)

	electron, err := s.repo.GetElectronFromRedis(param.Id)

	if err == redis.Nil {
		electron = &shared.User{
			Id: param.Id,
		}
	}

	if err != redis.Nil && err != nil {
		return nil, err
	}

	oldIndex := electron.LastKnownIndex

	electron.Coord = param.Coord
	electron.LastKnownIndex = newIndex

	err = s.repo.InsertElectronIntoRedis(electron)

	if err != nil {
		return nil, err
	}

	if oldIndex != newIndex {
		err = s.repo.RemoveElectronFromIndex(oldIndex, electron)
		if err != nil {
			return nil, err
		}

		err = s.repo.InsertElectronIntoIndex(newIndex, electron)
		if err != nil {
			return nil, err
		}

	}

	return electron, nil
}

func (s *Services) GetClosestElectrons(coord shared.Coord, steps int) []string {

	rings := geo.GetRingsFromOrigin(coord, steps)

	electronsIds := []string{}

	for _, index := range rings {
		ids, err := s.repo.GetElectronsInIndex(index)

		if err != nil {
			log.Printf("Failed to load electrons in electron_index %d", index)
			continue
		}

		if len(ids) > 0 {
			electronsIds = append(electronsIds, ids...)
		}
	}

	return electronsIds

}
