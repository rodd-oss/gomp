package main

import (
	"gomp_game/pkgs/gomp-v0/engine"
	"gomp_game/pkgs/gomp-v0/example/scenes"

	"time"
)

const tickRate = time.Second

func main() {
	e := engine.NewEngine(tickRate).SetDebug(false)

	v := e.LoadScene(scenes.VillageScene)
	v.ShouldRender = true

	e.RunClient()
}
