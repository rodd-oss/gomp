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
		return h.handleWsConnection(world, c.Response().Writer, c.Request())
	}
}

// serveWs handles websocket requests from the peer.
func (h *Hub) handleWsConnection(world *game.Game, w http.ResponseWriter, r *http.Request) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}

	unit := world.AddPlayer()
	defer world.RemovePlayer(unit)

	client := &Client{id: unit.Id, hub: h, conn: conn, send: make(chan []byte, 1024)}

	client.hub.register <- client
	defer func() { client.hub.unregister <- client }()

	event := &protos.Event{
		Type: protos.EventType_init,
		Data: &protos.Event_Init{
			Init: &protos.EventInit{
				PlayerId: unit.Id,
			},
		},
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

	event = &protos.Event{
		Type: protos.EventType_connect,
		Data: &protos.Event_Connect{
			Connect: &protos.EventConnect{Unit: unit},
		},
	}
	message, err = proto.Marshal(event)
	if err != nil {
		return err
	}
	h.broadcast <- message
	defer func() {
		event := &protos.Event{
			Type: protos.EventType_exit,
			Data: &protos.Event_Exit{
				Exit: &protos.EventExit{PlayerId: client.id},
			},
		}
		message, err := proto.Marshal(event)
		if err != nil {
			panic(err)
		}
		h.broadcast <- message
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	go client.writePump(&wg)
	go client.readPump(&wg, world)

	wg.Wait()

	return conn.Close()
}
