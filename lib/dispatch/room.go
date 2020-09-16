package dispatch

type RoomRequest struct {
	w    *WSConnection
	name string
}

type Room struct {
	name       string
	broadcast  chan []byte
	members    map[string]*WSConnection
	joinQueue  chan RoomRequest
	leaveQueue chan RoomRequest
	done       chan bool
}

func NewRoom(name string) *Room {
	room := &Room{
		name:       name,
		broadcast:  make(chan []byte),
		leaveQueue: make(chan RoomRequest),
		members:    make(map[string]*WSConnection),
		joinQueue:  make(chan RoomRequest),
		done:       make(chan bool),
	}

	go room.run()

	return room
}

func (room *Room) run() {
	for {
		select {

		case request := <-room.joinQueue:
			room.members[request.w.getIdBasedOnType()] = request.w
		case request := <-room.leaveQueue:
			if _, ok := room.members[request.w.getIdBasedOnType()]; ok {
				delete(room.members, request.w.getIdBasedOnType())
				close(request.w.Send)
			}
		case message := <-room.broadcast:
			for id, client := range room.members {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(room.members, id)
				}
			}

		case <-room.done:
			return
		}
	}
}

func (room *Room) close() {
	room.done <- true
}
