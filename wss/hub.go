package wss

type Hub struct {
	clients    map[*WSConnection]bool
	broadcast  chan []byte
	register   chan *WSConnection
	unregister chan *WSConnection
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *WSConnection),
		unregister: make(chan *WSConnection),
		clients:    make(map[*WSConnection]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
