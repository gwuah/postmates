package main

import (
	"log"
	"os"

	"github.com/gwuah/api/database/models"
	"github.com/gwuah/api/database/postgres"
	handler "github.com/gwuah/api/handlers"
	"github.com/gwuah/api/server"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	db, err := postgres.New(&postgres.Config{
		User:     dbUser,
		Password: dbPass,
		DBName:   dbName,
		SSLMode:  dbSSLMode,
		Host:     dbHost,
		Port:     dbPort,
	})

	if err != nil {
		log.Fatal("Failed To Connect To Postgresql database")
	}

	err = postgres.SetupDatabase(db, &models.Customer{}, &models.Delivery{}, &models.Electron{}, &models.Order{})

	if err != nil {
		log.Fatal("Failed To Setup Tables")
	}

	s := server.New()
	h := handler.New(db)

	routes := s.Group("/v1")
	h.Register(routes)

	server.Start(&s, &server.Config{
		Port: ":8080",
	})
}
