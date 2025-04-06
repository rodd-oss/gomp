/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.

===-===-===-===-===-===-===-===-===-===
Donations during this file development:
-===-===-===-===-===-===-===-===-===-===

<- Монтажер сука Donated 50 RUB

Thank you for your support!
*/

package stdsystems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
)

func NewSpriteSystem() SpriteSystem {
	return SpriteSystem{}
}

// SpriteSystem is a system that prepares Sprite to be rendered
type SpriteSystem struct {
	Positions     *stdcomponents.PositionComponentManager
	Scales        *stdcomponents.ScaleComponentManager
	Sprites       *stdcomponents.SpriteComponentManager
	RLTexturePros *stdcomponents.RLTextureProComponentManager
	Renderables   *stdcomponents.RenderableComponentManager
	RenderOrder   *stdcomponents.RenderOrderComponentManager
}

func (s *SpriteSystem) Init() {}
func (s *SpriteSystem) Run() {
	s.Sprites.EachEntity(func(entity ecs.Entity) bool {
		sprite := s.Sprites.Get(entity) //
		position := s.Positions.Get(entity)
		scale := s.Scales.Get(entity)

		renderable := s.Renderables.Get(entity)
		if renderable == nil {
			renderable = s.Renderables.Create(entity, stdcomponents.Renderable{
				Type: stdcomponents.SpriteRenderableType,
			})
		}

		renderOrder := s.RenderOrder.Get(entity)
		if renderOrder == nil {
			renderOrder = s.RenderOrder.Create(entity, stdcomponents.RenderOrder{})
		}

		tr := s.RLTexturePros.Get(entity)
		if tr == nil {
			s.RLTexturePros.Create(entity, stdcomponents.RLTexturePro{
				Texture: sprite.Texture, //
				Frame:   sprite.Frame,   //
				Origin: rl.Vector2{
					X: sprite.Origin.X * scale.XY.X,
					Y: sprite.Origin.Y * scale.XY.Y,
				},
				Dest: rl.Rectangle{X: position.XY.X, Y: position.XY.Y, Width: sprite.Frame.Width, Height: sprite.Frame.Height}, //
				Tint: sprite.Tint,
			})
		} else {
			tr.Dest.Width = sprite.Frame.Width
			tr.Dest.Height = sprite.Frame.Height
			tr.Tint = sprite.Tint
		}
		return true
	})
}
func (s *SpriteSystem) Destroy() {}
