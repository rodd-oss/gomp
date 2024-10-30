package main

import (
	"context"
	"gomp_game/pkgs/engine"
	"gomp_game/pkgs/scenes"
	"time"
)

const tickRate = time.Second / 4

func main() {
	e := engine.NewEngine(tickRate)

	e.Scenes = append(e.Scenes, scenes.VillageScene)

	ctx := context.Background()
	e.Run(ctx)
}
