package main

import (
	"github.com/gin-gonic/gin"
	engine "github.com/jilio/tomb_mates"
)

var world *engine.World

func init() {
	world = &engine.World{
		Replica: false,
		Units:   map[string]*engine.Unit{},
	}
}

func main() {
	go world.Evolve()

	hub := newHub()
	go hub.run()

	r := gin.New()
	r.GET("/ws", ginWsServe(hub, world))
	r.Run(":3000")
}

func ginWsServe(hub *Hub, world *engine.World) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		serveWs(hub, world, c.Writer, c.Request)
	})
}
