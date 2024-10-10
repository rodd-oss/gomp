package client

import (
	"context"
	"fmt"
	"log"
)

type Transport struct {
	Send     chan []byte
	Received chan []byte

	disconnect chan bool
	connect    chan bool

	ConnectHandler    func() error
	DisconnectHandler func() error

	IsConnected  bool
	IsConnecting bool
}

func newTransport(ctx context.Context) *Transport {
	t := &Transport{
		Send:         make(chan []byte),
		Received:     make(chan []byte),
		disconnect:   make(chan bool),
		connect:      make(chan bool),
		IsConnected:  false,
		IsConnecting: false,
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
			t.IsConnecting = true
			err := t.ConnectHandler()
			if err != nil {
				t.IsConnecting = false
				t.IsConnected = false
				log.Println("Transport connection error:", err)
				continue
			}
			t.IsConnecting = false
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

func (t *Transport) Connect() error {
	if t.IsConnected || t.IsConnecting {
		return fmt.Errorf("Already connected")
	}

	t.connect <- true

	return nil
}

func (t *Transport) Disconnect() error {
	if !t.IsConnected {
		return fmt.Errorf("Already disconnected")
	}

	t.disconnect <- true

	return nil
}
