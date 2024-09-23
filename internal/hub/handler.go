package hub

import (
	"log"
	"net/http"
	"tomb_mates/internal/engine"
	"tomb_mates/internal/protos"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"google.golang.org/protobuf/proto"
)

func (h *Hub) WsHandler(world *engine.World) echo.HandlerFunc {
	return func(c echo.Context) error {
		h.handleWsConnection(world, c.Response().Writer, c.Request())
		return nil
	}
}

// serveWs handles websocket requests from the peer.
func (h *Hub) handleWsConnection(world *engine.World, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	id := world.AddPlayer()
	client := &Client{id: id, hub: h, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	event := &protos.Event{
		Type: protos.Event_type_init,
		Data: &protos.Event_Init{
			Init: &protos.EventInit{
				PlayerId: id,
				Units:    world.Units,
			},
		},
	}
	world.Mx.Lock()
	message, err := proto.Marshal(event)
	world.Mx.Unlock()
	if err != nil {
		//todo: remove unit
		log.Println(err)
		return
	}
	conn.WriteMessage(websocket.BinaryMessage, message)

	unit := world.Units[id]
	event = &protos.Event{
		Type: protos.Event_type_connect,
		Data: &protos.Event_Connect{
			Connect: &protos.EventConnect{Unit: unit},
		},
	}
	world.Mx.Lock()
	message, err = proto.Marshal(event)
	world.Mx.Unlock()
	if err != nil {
		//todo: remove unit
		log.Println(err)
		return
	}

	h.broadcast <- message

	// Allow collection of memory referenced by the caller by doing all work
	// in new goroutines.
	go client.writePump()
	go client.readPump(world)
}
