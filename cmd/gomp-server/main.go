package main

import (
	"context"
	"gomp_game/pkgs/engine"
	"gomp_game/pkgs/example/scenes"

	"time"
)

const tickRate = time.Second / 1

func main() {
	e := engine.NewEngine(tickRate)
	e.SetDebug(true)

	e.LoadScene(scenes.VillageSceneName, scenes.VillageScene)

	e.Run(context.Background())
}
