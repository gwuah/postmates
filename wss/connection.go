package wss

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WSConnection struct {
	id             string
	hub            *Hub
	room           string
	conn           *websocket.Conn
	send           chan []byte
	processMessage func(msg []byte)
}

func (w *WSConnection) getIncomingMessages() {
	defer func() {
		log.Println("Unregistering  ...")
		w.hub.unregister <- w
		w.conn.Close()
	}()
	w.conn.SetReadLimit(maxMessageSize)
	w.conn.SetReadDeadline(time.Now().Add(pongWait))
	w.conn.SetPongHandler(func(string) error { w.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		go func() {
			w.processMessage(message)
		}()

		w.hub.broadcast <- message
	}
}

func (w *WSConnection) writeMessageToClient() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		w.conn.Close()
	}()
	for {
		select {
		case message, ok := <-w.send:
			w.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				w.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := w.conn.WriteMessage(websocket.TextMessage, message)

			if err != nil {
				log.Println("Failed to send message to client", err)
			}

		case <-ticker.C:
			w.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := w.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
