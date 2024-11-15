//go:build !server
// +build !server

/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	input "github.com/quasilyte/ebitengine-input"
)

func initInputSystem(cfg input.SystemConfig) input.System {
	sys := input.System{}
	sys.Init(cfg)
	return sys
}

var InputSystem = initInputSystem(input.SystemConfig{
	DevicesEnabled: input.AnyDevice,
})

func (g *Game) Ebiten() *ebitenGame {
	g.systems = append(g.systems, EbitenRenderSystem)
	EbitenRenderSystem.Init(g)

	e := new(ebitenGame)
	e.game = g
	e.inputSystem = &InputSystem

	tps := 1 / g.tickRate.Seconds()
	ebiten.SetTPS(int(tps))
	if g.Debug {
		log.Println("Initial TPS:", tps)
	}

	return e
}
