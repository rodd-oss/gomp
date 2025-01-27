package main

import (
	"gomp/examples/ebiten-ecs/scenes"
	systems2 "gomp/examples/ebiten-ecs/systems"
	"gomp/pkgs/gomp"

	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const tickRate = time.Second / 60

func main() {
	game := gomp.NewGame(tickRate)
	game.Debug = false

	game.LoadScene(scenes.VillageScene)

	// TODO: move BodySystem inside gomp engine such as render system
	game.RegisterSystems(
		gomp.BodySystem,
		systems2.IntentSystem,
		systems2.HeroMoveSystem,
	)

	e := game.Ebiten()

	e.RegisterInputHandlers()

	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Engine")

	err := ebiten.RunGame(e)
	if err != nil {
		panic(err)
	}
}
