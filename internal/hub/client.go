package hub

import (
	"log"
	"sync"
	"time"
	"tomb_mates/internal/game"
	"tomb_mates/internal/protos"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// // Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	id    string
	token string
	hub   *Hub
	conn  *websocket.Conn
	send  chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump(wg *sync.WaitGroup, world *game.Game) {
	defer func() {
		wg.Done()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return
		}

		event := &protos.Event{}
		err = proto.Unmarshal(message, event)
		if err != nil {
			return
		}

		world.RegisterEvent(event)
		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump(wg *sync.WaitGroup) {
	defer wg.Done()

	pingTicker := time.NewTicker(pingPeriod)
	defer pingTicker.Stop()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			err := c.conn.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				return
			}

		case <-pingTicker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
