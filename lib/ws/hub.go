package ws

import (
	"fmt"
	"sync"
)

type Hub struct {
	customers       map[string]*WSConnection
	couriers        map[string]*WSConnection
	broadcast       chan []byte
	Register        chan *WSConnection
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
		Register:   make(chan *WSConnection),
		unregister: make(chan *WSConnection),

		customers: make(map[string]*WSConnection),
		couriers:  make(map[string]*WSConnection),

		rooms:           make(map[string]*Room),
		createRoomQueue: make(chan RoomRequest),
		joinRoomQueue:   make(chan RoomRequest),
		leaveRoomQueue:  make(chan RoomRequest),

		gil: sync.Mutex{},
	}
}

func (h *Hub) GetSize(entities string) int {
	h.gil.Lock()
	defer h.gil.Unlock()

	if entities == "couriers" {
		return len(h.couriers)
	} else {
		return len(h.customers)
	}
}

func (h *Hub) GetCourier(id string) *WSConnection {
	// in future, we can refactor this so every entity has their own mutex
	h.gil.Lock()
	defer h.gil.Unlock()
	return h.couriers[id]
}

func (h *Hub) GetCustomer(id uint) *WSConnection {
	// in future, we can refactor this so every entity has their own mutex
	h.gil.Lock()
	defer h.gil.Unlock()
	return h.customers[fmt.Sprintf("%d", id)]
}

func (h *Hub) createRoom(name string) {
	if _, roomExists := h.rooms[name]; roomExists {
		return
	}

	h.rooms[name] = NewRoom(name)
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.Register:
			h.gil.Lock()

			if conn.Entity == "courier" {
				h.couriers[conn.Id] = conn
			} else {
				h.customers[conn.Id] = conn
			}

			h.gil.Unlock()

		case conn := <-h.unregister:
			h.gil.Lock()

			if conn.Entity == "courier" {
				if _, ok := h.couriers[conn.Id]; ok {
					delete(h.couriers, conn.Id)
					conn.Deactivate()
				}
			} else {
				if _, ok := h.customers[conn.Id]; ok {
					delete(h.customers, conn.Id)
					conn.Deactivate()
				}
			}

			h.gil.Unlock()

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
