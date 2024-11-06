/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package gomp

import (
	"image/color"

	"github.com/yohamta/donburi"
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
	})
}

func (c *ebitenRenderSystemController) Update(dt float64) {

}
