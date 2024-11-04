/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp_game/pkgs/gomp-v0/engine"
	"gomp_game/pkgs/gomp-v0/example/components"
	"gomp_game/pkgs/gomp-v0/example/entities"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/filter"
)

func RenderSystem() engine.RenderSystem {
	return engine.CreateRenderSystem(new(renderSystemController))
}

// renderSystemController is a system that updates the physics of a game
type renderSystemController struct {
	scene *engine.Scene
	world donburi.World
}

func (c *renderSystemController) Init(scene *engine.Scene) {
	c.scene = scene
	c.world = scene.World

	entities.PlayerRender.Each(c.world, func(e *donburi.Entry) {
		img := ebiten.NewImage(20, 20)

		img.Fill(color.RGBA{R: 255, G: 0, B: 0, A: 255})

		entities.PlayerRender.SetValue(e, components.RenderData{
			Sprite: img,
		})
	})
}

func (c *renderSystemController) Draw(screen *ebiten.Image, dt float64) {
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
