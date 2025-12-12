package websocket

import (
	"sync"
)

type Hub struct {
	rooms map[string]map[*Client]bool

	register   chan *Client
	unregister chan *Client
	broadcast  chan broadcastMsg

	mu sync.Mutex
}

type broadcastMsg struct {
	roomID  string
	message interface{}
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan broadcastMsg),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.rooms[client.RoomID]; !ok {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.rooms[client.RoomID]; ok {
				delete(h.rooms[client.RoomID], client)
				close(client.send)
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.Lock()
			if clients, ok := h.rooms[msg.roomID]; ok {
				for client := range clients {
					select {
					case client.send <- msg.message:
					default:
						close(client.send)
						delete(h.rooms[msg.roomID], client)
					}
				}
			}
			h.mu.Unlock()
		}
	}
}

func (h *Hub) Broadcast(roomID string, message interface{}) {
	h.broadcast <- broadcastMsg{roomID: roomID, message: message}
}
