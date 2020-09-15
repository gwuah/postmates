package wss

import (
	"log"
	"sync"
)

type Hub struct {
	clients         map[string]*WSConnection
	broadcast       chan []byte
	register        chan *WSConnection
	unregister      chan *WSConnection
	rooms           map[string]*Room
	createRoomQueue chan RoomRequest
	joinRoomQueue   chan RoomRequest
	leaveRoomQueue  chan RoomRequest
	gil             sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *WSConnection),
		unregister: make(chan *WSConnection),

		clients: make(map[string]*WSConnection),

		rooms:           make(map[string]*Room),
		createRoomQueue: make(chan RoomRequest),
		joinRoomQueue:   make(chan RoomRequest),
		leaveRoomQueue:  make(chan RoomRequest),

		gil: sync.Mutex{},
	}
}

func (h *Hub) createRoom(name string) {
	if _, roomExists := h.rooms[name]; roomExists {
		return
	}

	h.rooms[name] = NewRoom(name)
}

func (h *Hub) run() {
	for {
		select {
		case conn := <-h.register:
			log.Printf("Registering %s \n", conn.getIdBasedOnType())
			h.clients[conn.getIdBasedOnType()] = conn

		case conn := <-h.unregister:
			log.Printf("Unregistering %s\n", conn.getIdBasedOnType())

			if _, ok := h.clients[conn.getIdBasedOnType()]; ok {
				delete(h.clients, conn.getIdBasedOnType())
				close(conn.send)
			}

		case request := <-h.createRoomQueue:
			h.createRoom(request.name)

		case request := <-h.joinRoomQueue:
			room := h.rooms[request.name]
			room.joinQueue <- request

		case request := <-h.leaveRoomQueue:
			room := h.rooms[request.name]
			room.leaveQueue <- request
		}
	}
}
