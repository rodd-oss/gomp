//go:build !graphics
// +build !graphics

/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type ebitenGame struct {
	game         *Game
	dt           time.Duration
	lastUpdateAt time.Time
}

func (e *ebitenGame) Update() error {
	e.dt = time.Since(e.lastUpdateAt)

	e.game.Update(e.dt.Seconds())
	e.lastUpdateAt = time.Now()
	return nil
}

func (e *ebitenGame) Draw(screen *ebiten.Image) {
	e.game.Draw(screen)
}

func (e *ebitenGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return e.game.Layout(outsideWidth, outsideHeight)
}
