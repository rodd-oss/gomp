package main

import (
	"gomp_game/pkgs/example/scenes"
	"gomp_game/pkgs/gomp"

	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const tickRate = time.Second

func main() {
	game := gomp.NewGame(tickRate)
	game.Debug = true

	game.LoadScene(scenes.VillageScene)
	game.RegisterSystems(gomp.PhysicsSystem)

	e := game.Ebiten()
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Engine")

	err := ebiten.RunGame(e)
	if err != nil {
		panic(err)
	}
}
