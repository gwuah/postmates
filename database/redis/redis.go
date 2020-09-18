package redis

import (
	"github.com/go-redis/redis"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

func New(config *Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
}
