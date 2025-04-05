/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

<- HromRu Donated 1 500 RUB

Thank you for your support!
*/

package stdsystems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"time"
)

func NewRenderTexture2DRenderSystem() RenderTexture2DRenderSystem {
	return RenderTexture2DRenderSystem{}
}

type RenderTexture2DRenderSystem struct {
	EntityManager   *ecs.EntityManager
	RenderTexture2D *stdcomponents.RenderTexture2DComponentManager
}

func (s *RenderTexture2DRenderSystem) Init() {
}

func (s *RenderTexture2DRenderSystem) Run(dt time.Duration) bool {
	//if rl.IsKeyPressed(rl.KeyF12) {
	//	s.debug = !s.debug
	//}

	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.SetBlendMode(rl.BlendAdditive)
	s.RenderTexture2D.EachComponent(func(c *stdcomponents.RenderTexture2D) bool {
		rl.DrawTexturePro(c.Texture.Texture, rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  float32(c.Texture.Texture.Width),
			Height: -float32(c.Texture.Texture.Height),
		}, rl.Rectangle{
			X:      c.Position.X,
			Y:      c.Position.Y,
			Width:  c.Frame.Width,
			Height: c.Frame.Height,
		}, rl.Vector2{}, 0, rl.White)
		return true
	})

	rl.EndBlendMode()
	rl.EndDrawing()

	return false
}

func (s *RenderTexture2DRenderSystem) Destroy() {
}

type RenderTexture2DRenderInjector struct {
	EntityManager   *ecs.EntityManager
	RenderTexture2D *stdcomponents.RenderTexture2DComponentManager
}

func (s *RenderTexture2DRenderSystem) InjectWorld(injector *RenderTexture2DRenderInjector) {
	s.EntityManager = injector.EntityManager
	s.RenderTexture2D = injector.RenderTexture2D
}
