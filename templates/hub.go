package templates

import "realtime-chat/templates"

type Hub struct {
	clients    map[*templates.Client]bool
	broadcast  chan []byte
	register   chan *templates.Client
	unregister chan *templates.Client
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*templates.Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *templates.Client),
		unregister: make(chan *templates.Client),
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
