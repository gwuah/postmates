// This is basically our data layer.
package repository

import (
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type Repository struct {
	DB      *gorm.DB
	RedisDB *redis.Client
}

func New(db *gorm.DB, redisDB *redis.Client) *Repository {
	return &Repository{db, redisDB}
}
