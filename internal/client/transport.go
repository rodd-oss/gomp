package client

import (
	"context"
	"log"
)

type Transport struct {
	Send     chan []byte
	Received chan []byte

	disconnect chan bool
	connect    chan bool

	ConnectHandler    func() error
	DisconnectHandler func() error

	IsConnected bool
}

func newTransport(ctx context.Context) *Transport {
	t := &Transport{
		Send:        make(chan []byte),
		Received:    make(chan []byte),
		disconnect:  make(chan bool),
		connect:     make(chan bool),
		IsConnected: false,
	}

	go t.loop(ctx)

	return t
}

func (t *Transport) loop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		// case input := <-t.Send:
		// case output := <-t.Received:
		case <-t.connect:
			err := t.ConnectHandler()
			if err != nil {
				log.Println("Transport connection error:", err)
				t.IsConnected = false
				continue
			}
			t.IsConnected = true
		case <-t.disconnect:
			err := t.DisconnectHandler()
			if err != nil {
				log.Println("Transport disconnection error:", err)
				continue
			}

			t.IsConnected = false
		}
	}
}
