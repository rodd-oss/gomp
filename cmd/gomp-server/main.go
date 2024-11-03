package main

import (
	"context"
	"gomp_game/pkgs/engine"
	"gomp_game/pkgs/example/scenes"

	"time"
)

const tickRate = time.Second * 5

func main() {
	e := engine.NewEngine(tickRate).SetDebug(true)

	e.LoadScene(scenes.VillageScene)
	e.LoadScene(scenes.VillageScene)
	e.LoadScene(scenes.VillageScene)
	e.LoadScene(scenes.VillageScene)
	e.UnloadScene("") // this is a problem

	e.Run(context.Background())
}
