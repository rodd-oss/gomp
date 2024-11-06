/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

var EbitenRenderSystem = CreateSystem(new(ebitenRenderSystemController))

// ebitenRenderSystemController is a system that updates the physics of a game
type ebitenRenderSystemController struct {
	world donburi.World
}

func (c *ebitenRenderSystemController) Init(world donburi.World) {
	c.world = world

	RenderComponent.Each(c.world, func(e *donburi.Entry) {
		RenderComponent.Get(e).Sprite.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})
		log.Println("rendering sprite")
	})
}

func (c *ebitenRenderSystemController) Update(dt float64) {

}

func (c *ebitenRenderSystemController) Draw(screen *ebiten.Image, dt float64) {
	op := &ebiten.DrawImageOptions{}

	query := donburi.NewQuery(filter.Contains(PhysicsComponent, RenderComponent))

	query.Each(c.world, func(e *donburi.Entry) {
		render := RenderComponent.Get(e)
		physics := PhysicsComponent.Get(e)

		op.GeoM.Reset()
		op.GeoM.Translate(physics.Body.Position().X, physics.Body.Position().Y)

		screen.DrawImage(render.Sprite, op)
	})
}
