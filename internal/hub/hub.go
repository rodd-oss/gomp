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

const patchRate = time.Second / 20

func (h *Hub) run() {
	playersCount := -1

	playerCounterTicker := time.NewTicker(time.Second * 1)
	defer playerCounterTicker.Stop()

	patchTicker := time.NewTicker(patchRate)
	defer patchTicker.Stop()

	for {
		select {
		case <-playerCounterTicker.C:
			if playersCount != len(h.clients) {
				playersCount = len(h.clients)
				println("Players: ", playersCount)
			}
		case <-patchTicker.C:
			messages := make([]byte, 0)

			n := len(h.broadcast)
			if n == 0 {
				continue
			}

			for i := 0; i < n; i++ {
				message := <-h.broadcast
				messages = append(messages, message...)
			}

			for client := range h.clients {
				client.send <- messages
			}
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		}
	}
}
