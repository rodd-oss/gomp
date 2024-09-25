package hub

import (
	"net/http"
	"sync"
	"tomb_mates/internal/game"
	"tomb_mates/internal/protos"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
)

func (h *Hub) WsHandler(world *game.Game) echo.HandlerFunc {
	return func(c echo.Context) error {
		if len(world.Units) > 1000 {
			return c.String(http.StatusBadRequest, "Too many players")
		}
		return h.handleWsConnection(world, c.Response().Writer, c.Request())
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// serveWs handles websocket requests from the peer.
func (h *Hub) handleWsConnection(world *game.Game, w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	unit := world.AddPlayer()
	defer world.RemovePlayer(unit)

	client := &Client{id: unit.Id, hub: h, conn: conn, send: make(chan []byte, 512)}

	client.hub.register <- client
	defer func() { client.hub.unregister <- client }()

	event := &protos.Event{
		Type:     protos.EventType_init,
		PlayerId: unit.Id,
	}
	message, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, *world.UnitsSerialized)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go client.writePump(&wg)
	go client.readPump(&wg, world)

	wg.Wait()

	return conn.Close()
}
