package client

import (
	"context"
	"log"

	"github.com/coder/websocket"
)

// export transport

type wsTransport struct {
	ws        *websocket.Conn
	transport *Transport
}

func NewWsTransport(ctx context.Context, url string) *Transport {
	var wst = &wsTransport{}
	wst.transport = newTransport(ctx)

	wst.transport.ConnectHandler = func() error {
		err, w := wst.connect(ctx, url)
		if err != nil {
			return err
		}

		wst.ws = w
		return nil
	}

	wst.transport.DisconnectHandler = func() error {
		return wst.ws.Close(websocket.StatusNormalClosure, "Closed")
	}

	return wst.transport
}

func (wst *wsTransport) connect(ctx context.Context, url string) (error, *websocket.Conn) {
	w, _, err := websocket.Dial(ctx, url, nil)
	if err != nil {
		return err, nil
	}

	// Write messages
	go func(c *websocket.Conn) {
		defer c.CloseNow()

		c.SetReadLimit(-1)

		for {
			select {
			case <-ctx.Done():
				return

			case message := <-wst.transport.Send:
				err = c.Write(ctx, websocket.MessageBinary, message)
				if err != nil {
					log.Println("Error writing message:", err)
					return
				}
			}
		}
	}(w)

	// Read messages
	go func(c *websocket.Conn) {
		defer c.CloseNow()

		c.SetReadLimit(-1)

		for {
			var _, message, err = c.Read(ctx)
			if err != nil {
				log.Println("Error reading message:", err)
				return
			}

			wst.transport.Received <- message
		}
	}(w)

	return nil, w
}
