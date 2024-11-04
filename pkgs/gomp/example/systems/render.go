/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp_game/pkgs/gomp/ecs"
	"gomp_game/pkgs/gomp/example/components"
	"gomp_game/pkgs/gomp/example/entities"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

func EbitenRenderSystem() ecs.System {
	// TODO Implement CreateRenderSystem
	// return gomp.CreateSystem(new(renderSystemController))
	return ecs.System{}
}

// ebitenRenderSystemController is a system that updates the physics of a game
type ebitenRenderSystemController struct {
	world donburi.World
}

func (c *ebitenRenderSystemController) Init(world donburi.World) {
	c.world = world

	entities.PlayerRender.Each(c.world, func(e *donburi.Entry) {
		img := ebiten.NewImage(20, 20)

		img.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})

		entities.PlayerRender.SetValue(e, components.RenderData{
			Sprite: img,
		})
	})
}

func (c *ebitenRenderSystemController) Draw(screen *ebiten.Image, dt float64) {
	op := &ebiten.DrawImageOptions{}

	query := donburi.NewQuery(filter.Contains(entities.PlayerPhysics, entities.PlayerRender))

	query.Each(c.world, func(e *donburi.Entry) {
		render := entities.PlayerRender.Get(e)
		physics := entities.PlayerPhysics.Get(e)

		op.GeoM.Reset()
		op.GeoM.Translate(physics.Body.Position().X, physics.Body.Position().Y)

		screen.DrawImage(render.Sprite, op)
	})
}
