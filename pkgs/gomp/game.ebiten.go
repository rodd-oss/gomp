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
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
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

func (g *Game) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}

	query := donburi.NewQuery(filter.Contains(BodyComponent, RenderComponent))

	query.Each(g.world, func(e *donburi.Entry) {
		render := RenderComponent.Get(e)

		if render == nil {
			log.Fatalln("RenderComponent is nil")
		}

		body := BodyComponent.Get(e)

		if body == nil {
			log.Fatalln("BodyComponent is nil")
		}

		op.GeoM.Reset()
		op.GeoM.Translate(body.Position().X, body.Position().Y)

		screen.DrawImage(render.Sprite, op)
	})
}

func (c *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
