package repository

import (
	"encoding/json"
	"fmt"

	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/shared"
	"github.com/uber/h3-go"
)

func (r *Repository) UpdateElectron(id uint, data map[string]interface{}) (*models.Electron, error) {
	electron := models.Electron{}

	if err := r.DB.Model(&electron).Where("id = ?", id).Updates(data).Error; err != nil {
		return nil, err
	}

	return &electron, nil
}

func (r *Repository) GetElectronFromRedis(id string) (*shared.User, error) {
	var user shared.User
	key := fmt.Sprintf("electron_%s", id)

	result, err := r.RedisDB.Get(key).Result()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(result), &user)

	if err != nil {
		return nil, err
	}

	return &user, nil

}

func (r *Repository) InsertElectronIntoRedis(user *shared.User) error {
	stringifiedUser, err := json.Marshal(user)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("electron_%s", user.Id)
	_, err = r.RedisDB.Set(key, stringifiedUser, 0).Result()

	if err != nil {
		return err
	}
	return nil

}

func (r *Repository) RemoveElectronFromIndex(index h3.H3Index, user *shared.User) error {
	key := fmt.Sprintf("electron_index_%d", index)
	_, err := r.RedisDB.LRem(key, 0, user.Id).Result()
	if err != nil {
		return err
	}
	return nil

}

func (r *Repository) InsertElectronIntoIndex(index h3.H3Index, user *shared.User) error {
	key := fmt.Sprintf("electron_index_%d", index)
	_, err := r.RedisDB.LPush(key, user.Id).Result()
	if err != nil {
		return err
	}
	return nil

}

func (r *Repository) GetElectronsInIndex(index h3.H3Index) ([]string, error) {
	key := fmt.Sprintf("electron_index_%d", index)
	electronsIds, err := r.RedisDB.LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return electronsIds, nil

}
