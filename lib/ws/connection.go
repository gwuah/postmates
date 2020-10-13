package ws

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024
)

type WSConnection struct {
	Id                    string
	Hub                   *Hub
	Room                  string
	Conn                  *websocket.Conn
	Send                  chan []byte
	ProcessMessage        func(msg []byte, ws *WSConnection)
	Entity                string
	IsActive              bool
	DeliveryAcceptanceAck chan bool
}

func (w *WSConnection) Deactivate() {
	close(w.Send)
	w.IsActive = false
}

func (w *WSConnection) ReadPump() {
	defer func() {
		log.Println("Unregistering", w.Id)
		w.Hub.unregister <- w
		w.Conn.Close()
	}()
	w.Conn.SetReadLimit(maxMessageSize)
	w.Conn.SetReadDeadline(time.Now().Add(pongWait))
	w.Conn.SetPongHandler(func(string) error { w.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := w.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		go func() {
			w.ProcessMessage(message, w)
		}()

	}
}

func (w *WSConnection) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		w.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-w.Send:
			w.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The Hub closed the channel.
				w.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := w.Conn.WriteMessage(websocket.TextMessage, message)

			if err != nil {
				log.Println("failed to Send message to client", err)
			}

		case <-ticker.C:
			w.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := w.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (w *WSConnection) GetIdBasedOnType() string {
	if w.Entity == "courier" {
		return fmt.Sprintf("courier_%s", w.Id)
	} else {
		return fmt.Sprintf("customer_%s", w.Id)
	}
}

func (w *WSConnection) JoinRoom(name string) {
	w.Hub.joinRoomQueue <- RoomRequest{name: name, w: w}
}

func (w *WSConnection) LeaveRoom(name string) {
	w.Hub.leaveRoomQueue <- RoomRequest{name: name, w: w}
}

func (w *WSConnection) SendMessage(message []byte) {
	if w.IsActive {
		w.Send <- message
	} else {
		log.Println("Can't send message to closed socket conn", w.Id, w.Entity)
	}
}

func (w *WSConnection) AckDeliveryAcceptance(status bool) {
	if w.IsActive {
		w.DeliveryAcceptanceAck <- status
	} else {
		log.Println("Can't send message to closed socket conn", w.Id, w.Entity)
	}
}
