/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package stdsystems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
)

func NewSpriteMatrixSystem() SpriteMatrixSystem {
	return SpriteMatrixSystem{}
}

// SpriteMatrixSystem is a system that prepares SpriteSheet to be rendered
type SpriteMatrixSystem struct {
	SpriteMatrixes  *stdcomponents.SpriteMatrixComponentManager
	RLTexturePros   *stdcomponents.RLTextureProComponentManager
	AnimationStates *stdcomponents.AnimationStateComponentManager
	Positions       *stdcomponents.PositionComponentManager
}

func (s *SpriteMatrixSystem) Init() {}
func (s *SpriteMatrixSystem) Run() {
	s.SpriteMatrixes.EachEntity(func(entity ecs.Entity) bool {
		spriteMatrix := s.SpriteMatrixes.Get(entity) //
		position := s.Positions.GetUnsafe(entity)
		animationState := s.AnimationStates.GetUnsafe(entity) //

		frame := spriteMatrix.Animations[*animationState].Frame

		tr := s.RLTexturePros.GetUnsafe(entity)
		if tr == nil {
			s.RLTexturePros.Create(entity, stdcomponents.RLTexturePro{
				Texture: spriteMatrix.Texture, //
				Frame:   frame,                //
				Origin:  spriteMatrix.Origin,
				Dest:    rl.Rectangle{X: position.XY.X, Y: position.XY.Y, Width: frame.Width, Height: frame.Height}, //
			})
		} else {
			// Run spriteRender

			tr.Frame = frame
		}
		return true
	})
}
func (s *SpriteMatrixSystem) Destroy() {}
