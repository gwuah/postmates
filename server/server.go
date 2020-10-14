package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/electra-systems/core-api/middleware"
	"github.com/electra-systems/core-api/utils/validator"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Config struct {
	Port  string
	Debug bool
}

type Server struct {
	*gin.Engine
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "OK",
	})
}

func New() Server {
	binding.Validator = new(validator.DefaultValidator)
	server := gin.Default()
	server.Use(middleware.CORS())
	server.GET("/", healthCheck)
	return Server{server}
}

func Start(e *Server, cfg *Config) {

	s := &http.Server{
		Addr:    cfg.Port,
		Handler: e.Engine,
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		if err := s.Close(); err != nil {
			log.Println("failed To ShutDown Server", err)
		}
		log.Println("Shut Down Server")
	}()

	if err := s.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server Closed After Interruption")
		} else {
			log.Println("Unexpected Server Shutdown")
		}
	}
}
