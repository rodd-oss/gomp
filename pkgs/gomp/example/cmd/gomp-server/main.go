package main

import (
	"context"
	"gomp_game/pkgs/gomp"
	"gomp_game/pkgs/gomp/example/scenes"

	"time"
)

const tickRate = time.Second

func main() {
	e := gomp.NewGame(tickRate)
	e.Debug = true

	e.LoadScene(scenes.VillageScene)

	e.Run(context.Background())
}
