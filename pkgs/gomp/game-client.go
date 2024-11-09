//go:build !graphics
// +build !graphics

/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Ebiten() *ebitenGame {
	g.systems = append(g.systems, EbitenRenderSystem)
	EbitenRenderSystem.Init(g.world)

	e := new(ebitenGame)
	e.game = g

	tps := 1 / g.tickRate.Seconds()
	log.Println(tps)
	ebiten.SetTPS(int(tps))
	if g.Debug {
		log.Println("Initial TPS:", tps)
	}

	return e
}
