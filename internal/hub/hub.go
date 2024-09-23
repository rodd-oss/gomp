package hub

import (
	"time"
)

// Hub maintains the set of active clients and broadcasts messages
// to the clients.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func New() *Hub {
	h := &Hub{
		broadcast:  make(chan []byte, 1),
		register:   make(chan *Client, 1),
		unregister: make(chan *Client, 1),
		clients:    make(map[*Client]bool),
	}

	go h.run()

	return h
}

func (h *Hub) run() {
	ticket := time.NewTicker(time.Second * 1)
	playersCount := -1
	for {
		select {
		case <-ticket.C:
			if playersCount != len(h.clients) {
				playersCount = len(h.clients)
				println("Players: ", playersCount)
			}
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
