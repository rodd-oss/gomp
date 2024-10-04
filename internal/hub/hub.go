package hub

import (
	"time"
	"tomb_mates/internal/game"
)

// Hub maintains the set of active clients and broadcasts messages
// to the clients.
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func New(game *game.Game) *Hub {
	h := &Hub{
		broadcast:  make(chan []byte, 1),
		register:   make(chan *Client, 32),
		unregister: make(chan *Client, 32),
		clients:    make(map[*Client]bool),
	}

	go h.run(game)

	return h
}

func (h *Hub) run(game *game.Game) {
	playersCount := -1

	playerCounterTicker := time.NewTicker(time.Second * 1)
	defer playerCounterTicker.Stop()

	for {
		select {
		case <-playerCounterTicker.C:
			if playersCount != len(h.clients) {
				playersCount = len(h.clients)
				println("Players: ", playersCount)
			}
		case message := <-game.NetworkManager.Broadcast:
			for client := range h.clients {
				if len(client.send) != cap(client.send) {
					client.send <- message
				}
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				if len(client.send) != cap(client.send) {
					client.send <- message
				}
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
