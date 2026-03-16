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
	"github.com/negrel/assert"
	"gomp/pkg/core"
	"gomp/pkg/ecs"
	"gomp/pkg/worker"
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
	RenderOrder   *stdcomponents.RenderOrderComponentManager

	numWorkers       int
	accRenderOrder   [][]ecs.Entity
	accRLTexturePros [][]ecs.Entity
	Engine           *core.Engine
}

func (s *SpriteSystem) Init() {
	s.numWorkers = s.Engine.Pool().NumWorkers()
	s.accRenderOrder = make([][]ecs.Entity, s.numWorkers)
	s.accRLTexturePros = make([][]ecs.Entity, s.numWorkers)
}
func (s *SpriteSystem) Run() {
	for i := range s.accRenderOrder {
		s.accRenderOrder[i] = s.accRenderOrder[i][:0]
	}
	for i := range s.accRLTexturePros {
		s.accRLTexturePros[i] = s.accRLTexturePros[i][:0]
	}
	s.Sprites.ProcessEntities(func(entity ecs.Entity, workerId worker.WorkerId) {
		if !s.RenderOrder.Has(entity) {
			s.accRenderOrder[workerId] = append(s.accRenderOrder[workerId], entity)
		}
		if !s.RLTexturePros.Has(entity) {
			s.accRLTexturePros[workerId] = append(s.accRLTexturePros[workerId], entity)
		}
	})
	for a := range s.accRenderOrder {
		for _, entity := range s.accRenderOrder[a] {
			s.RenderOrder.Create(entity, stdcomponents.RenderOrder{})
		}
	}
	for a := range s.accRLTexturePros {
		for _, entity := range s.accRLTexturePros[a] {
			s.RLTexturePros.Create(entity, stdcomponents.RLTexturePro{})
		}
	}

	s.Sprites.ProcessEntities(s.updateTextureRender)
}
func (s *SpriteSystem) Destroy() {}
func (s *SpriteSystem) updateTextureRender(entity ecs.Entity, workerId worker.WorkerId) {
	sprite := s.Sprites.GetUnsafe(entity)
	assert.NotNil(sprite)

	position := s.Positions.GetUnsafe(entity)
	assert.NotNil(position)

	scale := s.Scales.GetUnsafe(entity)
	assert.NotNil(scale)

	renderOrder := s.RenderOrder.GetUnsafe(entity)
	assert.NotNil(renderOrder)

	tr := s.RLTexturePros.GetUnsafe(entity)
	assert.NotNil(tr)

	tr.Texture = sprite.Texture
	tr.Frame = sprite.Frame
	tr.Origin = rl.Vector2{
		X: sprite.Origin.X * scale.XY.X,
		Y: sprite.Origin.Y * scale.XY.Y,
	}
	tr.Dest.X = position.XY.X
	tr.Dest.Y = position.XY.Y
	tr.Dest.Width = sprite.Frame.Width
	tr.Dest.Height = sprite.Frame.Height
	tr.Tint = sprite.Tint
}
