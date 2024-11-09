package main

import (
	"gomp_game/pkgs/example/scenes"
	"gomp_game/pkgs/gomp"

	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const tickRate = time.Second / 1

func main() {
	game := gomp.NewGame(tickRate)
	game.Debug = true

	game.LoadScene(scenes.VillageScene)

	// TODO: move BodySystem inside gomp engine such as render system
	game.RegisterSystems(
		gomp.BodySystem,
	)

	e := game.Ebiten()
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Engine")

	err := ebiten.RunGame(e)
	if err != nil {
		panic(err)
	}
}
