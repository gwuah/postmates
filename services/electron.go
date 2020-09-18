package services

import (
	"github.com/go-redis/redis"
	"github.com/gwuah/api/shared"
	"github.com/gwuah/api/utils/geo"
)

func (s *Services) IndexElectronLocation(param shared.UserLocationUpdate) (*shared.User, error) {
	newIndex := geo.CoordToIndex(param)

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
