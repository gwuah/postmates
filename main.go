package main

import (
	handler "github.com/gwuah/api/handlers"
	"github.com/gwuah/api/models"
	"github.com/gwuah/api/postgres"
	"github.com/gwuah/api/server"
)

func main() {

	db, err := postgres.New(&postgres.Config{
		User:     "postgres",
		Password: "mtn_loyalty",
		DBName:   "electra_dev",
		SSLMode:  "disable",
		Host:     "127.0.0.1",
	})

	if err != nil {
		panic("Failed To Connect To PostGresSQL")
	}

	err = postgres.SetupDatabase(db, &models.Customer{}, &models.Delivery{}, &models.Electron{}, &models.Order{})

	if err != nil {
		panic("Failed To Setup Tables")
	}

	s := server.New()
	e := s.Group("/v1")
	h := handler.New(db)
	h.Register(e)

	server.Start(&s, &server.Config{
		Port: ":8080",
	})
}
