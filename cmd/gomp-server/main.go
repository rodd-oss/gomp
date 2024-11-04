package main

import (
	"context"
	"gomp_game/pkgs/engine"
	"gomp_game/pkgs/example/scenes"

	"time"
)

const tickRate = time.Second

func main() {
	e := engine.NewEngine(tickRate).SetDebug(true).SetDebugDraw(false)

	e.LoadScene(scenes.VillageScene)

	e.Run(context.Background())
}
