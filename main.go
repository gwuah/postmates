package main

import (
	"crypto/sha1"
	"log"
	"os"

	"github.com/gwuah/api/database"
	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/database/postgres"
	"github.com/gwuah/api/handler"
	"github.com/gwuah/api/server"
	"github.com/gwuah/api/utils/jwt"
	"github.com/gwuah/api/utils/secure"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := postgres.New(&postgres.Config{
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
	})

	if err != nil {
		log.Fatal("Failed To Connect To Postgresql database")
	}

	err = postgres.SetupDatabase(db, &models.Customer{}, &models.Delivery{}, &models.Electron{}, &models.Order{})

	if err != nil {
		log.Fatal("Failed To Setup Tables")
	}

	database.RunSeeds(db, []database.SeedFn{
		database.SeedProducts,
		database.SeedElectrons,
		database.SeedCustomers,
	})

	sec := secure.New(1, sha1.New())

	jwt, err := jwt.New("HS256", os.Getenv("JWT_SECRET"), 15, 64)
	if err != nil {
		log.Fatal(err)
	}

	s := server.New()
	h := handler.New(db, jwt, sec)

	routes := s.Group("/v1")
	h.Register(routes)

	server.Start(&s, &server.Config{
		Port: ":8080",
	})
}
