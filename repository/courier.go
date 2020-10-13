package repository

import (
	"encoding/json"
	"fmt"

	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/shared"
	"github.com/uber/h3-go"
	"gorm.io/gorm/clause"
)

func (r *Repository) FindCourier(id uint) (*models.Courier, error) {
	courier := models.Courier{}

	if err := r.DB.Preload(clause.Associations).First(&courier, id).Error; err != nil {
		return nil, err
	}

	return &courier, nil
}

func (r *Repository) UpdateCourier(id uint, data map[string]interface{}) (*models.Courier, error) {
	courier := models.Courier{}

	if err := r.DB.Model(&courier).Where("id = ?", id).Updates(data).Error; err != nil {
		return nil, err
	}

	return &courier, nil
}

func (r *Repository) GetCourierFromRedis(id string) (*shared.User, error) {
	var user shared.User
	key := fmt.Sprintf("courier_%s", id)

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

func (r *Repository) InsertCourierIntoRedis(user *shared.User) error {
	stringifiedUser, err := json.Marshal(user)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("courier_%s", user.Id)
	_, err = r.RedisDB.Set(key, stringifiedUser, 0).Result()

	if err != nil {
		return err
	}
	return nil

}

func (r *Repository) RemoveCourierFromIndex(index h3.H3Index, user *shared.User) error {
	key := fmt.Sprintf("courier_index_%d", index)
	_, err := r.RedisDB.LRem(key, 0, user.Id).Result()
	if err != nil {
		return err
	}
	return nil

}

func (r *Repository) InsertCourierIntoIndex(index h3.H3Index, user *shared.User) error {
	key := fmt.Sprintf("courier_index_%d", index)
	_, err := r.RedisDB.LPush(key, user.Id).Result()
	if err != nil {
		return err
	}
	return nil

}

func (r *Repository) GetCouriersInIndex(index h3.H3Index) ([]string, error) {
	key := fmt.Sprintf("courier_index_%d", index)
	couriersIds, err := r.RedisDB.LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}
	return couriersIds, nil

}

func (r *Repository) GetAllCouriers(ids []string) ([]*shared.User, error) {
	couriers := []*shared.User{}

	for _, id := range ids {
		courier, _ := r.GetCourierFromRedis(id)
		couriers = append(couriers, courier)
	}

	return couriers, nil
}
