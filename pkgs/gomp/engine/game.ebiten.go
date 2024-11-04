//go:build !graphics
// +build !graphics

/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package engine

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Ebiten() *ebitenGame {
	// g.systems = append(g.systems, systems.RenderSystem())

	e := new(ebitenGame)
	e.game = g

	tps := int(1000 / e.game.tickRate.Milliseconds())
	ebiten.SetTPS(tps)
	if g.Debug {
		log.Println("Initial TPS:", tps)
	}

	return e
}

func (g *Game) Draw(screen *ebiten.Image) {
	// renderSystem.Draw()
}

func (c *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
