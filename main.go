package main

import (
	"github.com/gwuah/api/postgres"
	"github.com/gwuah/api/server"
)

func main() {

	_, err := postgres.New(&postgres.Config{
		User:    "user",
		DBName:  "electra_dev",
		SSLMode: "disable",
		Host:    "127.0.0.1",
	})

	if err != nil {
		panic("Failed To Connect To PostGresSQL")
	}

	s := server.New()
	c := &server.Config{
		Port: ":8080",
	}
	server.Start(&s, c)
}
