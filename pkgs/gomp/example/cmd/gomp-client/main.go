package main

import (
	"gomp_game/pkgs/gomp"
	"gomp_game/pkgs/gomp/example/scenes"

	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const tickRate = time.Second

func main() {
	game := gomp.NewGame(tickRate)
	game.Debug = true

	game.LoadScene(scenes.VillageScene)

	e := game.Ebiten()
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Engine")

	if err := ebiten.RunGame(e); err != nil {
		panic(err)
	}
}
