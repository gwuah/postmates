package main

import (
	"github.com/gwuah/api/server"
)

func main() {
	s := server.New()
	c := &server.Config{
		Port: ":8080",
	}
	server.Start(&s, c)
}
