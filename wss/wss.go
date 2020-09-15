package wss

import (
	"log"

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

func (wss *WSS) HandleWebsocketConnection(c *gin.Context) {
	id := c.Param("id")
	log.Println("Connection Recieved from", id)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		log.Println("Failed to setup websocket conn ..", err)
		return
	}

	wsConnection := &WSConnection{hub: wss.hub, send: make(chan []byte), conn: conn, id: id}

	go wsConnection.getIncomingMessages()
	go wsConnection.writeMessageToClient()

}
