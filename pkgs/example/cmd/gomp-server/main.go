package main

import (
	"context"
	"gomp_game/pkgs/example/scenes"
	"gomp_game/pkgs/gomp"

	"time"
)

const tickRate = time.Second

func main() {
	e := gomp.NewGame(tickRate)
	e.Debug = true

	e.LoadScene(scenes.VillageScene)

	e.Run(context.Background())
}
