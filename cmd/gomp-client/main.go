package main

import (
	"gomp_game/pkgs/engine"
	"gomp_game/pkgs/example/scenes"

	"time"
)

const tickRate = time.Second

func main() {
	e := engine.NewEngine(tickRate).SetDebug(false)

	v := e.LoadScene(scenes.VillageScene)
	v.ShouldRender = true

	e.RunClient()
}
