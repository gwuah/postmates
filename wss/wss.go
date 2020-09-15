package wss

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type WSS struct {
	hub *Hub
}

func New() *WSS {
	hub := NewHub()
	go hub.run()
	return &WSS{hub}
}

func (wss *WSS) HandleWebsocketConnection(entity string) func(c *gin.Context) {
	return func(c *gin.Context) {
		id := c.Param("id")
		log.Println("Connection Recieved from", id)
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

		if err != nil {
			log.Println("Failed to setup websocket conn ..", err)
			return
		}

		wsConnection := &WSConnection{
			hub:    wss.hub,
			send:   make(chan []byte),
			conn:   conn,
			id:     id,
			entity: entity,
		}

		wss.hub.register <- wsConnection

		go wsConnection.getIncomingMessages()
		go wsConnection.writeMessageToClient()

		go func() {
			ticker := time.NewTicker(1 * time.Second)
			done := make(chan bool)

			go func() {
				for {
					select {
					case value := <-ticker.C:
						wsConnection.sendMessage([]byte(fmt.Sprintf("%s %s", "message", value)))
					case <-done:
						return
					}
				}
			}()

			time.Sleep(6 * time.Second)
			close(done)
		}()
	}
}
