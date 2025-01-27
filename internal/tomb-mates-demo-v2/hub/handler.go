/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package hub

import (
	"gomp/internal/tomb-mates-demo-v2/game"
	"gomp/internal/tomb-mates-demo-v2/protos"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
)

func (h *Hub) WsHandler(world *game.Game) echo.HandlerFunc {
	return func(c echo.Context) error {
		if len(world.Entities.Units) >= int(world.MaxPlayers) {
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

	id := world.GeneratePlayerId()
	world.CreatePlayer(id, nil)
	defer world.RemovePlayer(id)

	client := &Client{id: id, hub: h, conn: conn, send: make(chan []byte, 512)}

	client.hub.register <- client
	defer func() { client.hub.unregister <- client }()

	event := &protos.Event{
		Type:     protos.EventType_init,
		PlayerId: id,
	}
	message, err := proto.Marshal(event)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, message)
	if err != nil {
		return err
	}

	err = conn.WriteMessage(websocket.BinaryMessage, *world.StateSerialized)
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
