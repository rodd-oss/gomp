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

	e.LoadScene(scenes.VillageSceneName, scenes.VillageScene)
	e.LoadScene(scenes.DungeonSceneName, scenes.DungeonScene)

	go after(e)

	ctx := context.Background()
	e.Run(ctx)
}

func after(e *engine.Engine) {
	t := time.After(time.Second * 5)
	<-t

	e.UnloadAllScenes()
}
