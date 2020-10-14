package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/electra-systems/core-api/database"
	"github.com/electra-systems/core-api/database/models"
	"github.com/electra-systems/core-api/database/postgres"
	"github.com/electra-systems/core-api/database/redis"
	"github.com/electra-systems/core-api/handler"
	"github.com/electra-systems/core-api/server"
	"github.com/electra-systems/core-api/utils/jwt"
	"github.com/electra-systems/core-api/utils/secure"
	"github.com/joho/godotenv"
)

func main() {
	ENV := os.Getenv("ENV")
	fmt.Println("ENVVVV", ENV)

	if (ENV != "staging") || (ENV != "production") {
		log.Println("loading env variables")
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file", err)
		}
	}

	db, err := postgres.New(&postgres.Config{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		DBurl:    os.Getenv("DATABASE_URL"),
	})

	if err != nil {
		log.Fatal("failed To Connect To Postgresql database", err)
	}

	err = postgres.SetupDatabase(db,
		&models.Customer{},
		&models.Delivery{},
		&models.Courier{},
		&models.Order{},
		&models.Vehicle{},
		&models.TripPoint{},
	)

	if err != nil {
		log.Fatal("failed To Setup Tables")
	}

	database.RunSeeds(db, []database.SeedFn{
		database.SeedProducts,
		database.SeedCouriers,
		database.SeedCustomers,
		database.SeedVehicles,
	})

	sec := secure.New(1, sha1.New())

	jwt, err := jwt.New("HS256", os.Getenv("JWT_SECRET"), 15, 64)

	if err != nil {
		log.Fatal(err)
	}

	REDIS_DB, err := strconv.Atoi(os.Getenv("REDIS_DB"))

	if err != nil {
		log.Fatal(err)
	}

	redisDB := redis.New(&redis.Config{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       REDIS_DB,
	})

	if err != nil {
		log.Fatal("failed to initialize mapbox", err)
	}

	s := server.New()
	h := handler.New(db, jwt, sec, redisDB)

	routes := s.Group("/v1")
	h.Register(routes)

	server.Start(&s, &server.Config{
		Port: ":8080",
	})
}
