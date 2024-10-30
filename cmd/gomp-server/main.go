package main

import (
	"context"
	"gomp_game/pkgs/engine"
	"gomp_game/pkgs/scenes"
	"time"
)

const tickRate = time.Second / 1

func main() {
	e := engine.NewEngine(tickRate)

	e.RegisterScene("Village", scenes.VillageScene)
	e.RegisterScene("Dungeon", scenes.DungeonScene)

	ctx := context.Background()
	e.Run(ctx)
}
