package main

import (
	"context"
	"gomp/examples/ebiten-ecs/scenes"
	"gomp/pkg/gomp"
)

func main() {
	e := gomp.NewGame(tickRate)
	e.Debug = true

	e.LoadScene(scenes.VillageScene)

	e.Run(context.Background())
}
