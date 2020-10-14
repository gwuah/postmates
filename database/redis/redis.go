package redis

import (
	"net/url"
	"os"

	"github.com/go-redis/redis"
)

type Config struct {
	Addr     string
	Password string
	DB       int
	DBurl    string
}

func New(config *Config) *redis.Client {
	if os.Getenv("ENV") == "staging" || os.Getenv("ENV") == "production" {
		parsedURL, _ := url.Parse(config.DBurl)
		password, _ := parsedURL.User.Password()
		return redis.NewClient(&redis.Options{
			Addr:     parsedURL.Host,
			Password: password,
		})
	}

	return redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
}
