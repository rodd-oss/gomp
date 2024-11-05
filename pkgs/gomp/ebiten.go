//go:build !graphics
// +build !graphics

/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

import "github.com/hajimehoshi/ebiten/v2"

type ebitenGame struct {
	game *Game
}

func (e *ebitenGame) Update() error {
	tps := ebiten.ActualTPS()
	e.game.Update(tps)

	return nil
}

func (e *ebitenGame) Draw(screen *ebiten.Image) {
	e.game.Draw(screen)
}

func (e *ebitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return e.game.Layout(outsideWidth, outsideHeight)
}
