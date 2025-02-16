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
	"time"
)

func NewTextureRenderMatrixSystem(
	spriteMatrixes *stdcomponents.SpriteMatrixComponentManager,
	textureRenders *stdcomponents.TextureRenderComponentManager,
	animationStates *stdcomponents.AnimationStateComponentManager,
) *TextureRenderMatrixSystem {
	return &TextureRenderMatrixSystem{
		spriteMatrixes:  spriteMatrixes,
		textureRenders:  textureRenders,
		animationStates: animationStates,
	}
}

// TextureRenderMatrixSystem is a system that prepares SpriteSheet to be rendered
type TextureRenderMatrixSystem struct {
	spriteMatrixes  *stdcomponents.SpriteMatrixComponentManager
	textureRenders  *stdcomponents.TextureRenderComponentManager
	animationStates *stdcomponents.AnimationStateComponentManager
}

func (s *TextureRenderMatrixSystem) Init() {}
func (s *TextureRenderMatrixSystem) Run(dt time.Duration) {
	// Run sprites and spriteRenders
	s.spriteMatrixes.AllParallel(func(entity ecs.Entity, spriteMatrix *stdcomponents.SpriteMatrix) bool {
		if spriteMatrix == nil {
			return true
		}

		animationState := s.animationStates.Get(entity)
		if animationState == nil {
			return true
		}

		currentAnimationFrame := spriteMatrix.Animations[*animationState].Frame

		tr := s.textureRenders.Get(entity)
		if tr == nil {
			// Create new spriteRender
			newRender := stdcomponents.TextureRender{
				Texture: spriteMatrix.Texture,
				Origin:  spriteMatrix.Origin,
				Frame:   currentAnimationFrame,
				Dest: rl.Rectangle{
					Width:  currentAnimationFrame.Width,
					Height: currentAnimationFrame.Height,
				},
			}

			s.textureRenders.Create(entity, newRender)
		} else {
			// Run spriteRender
			tr.Texture = spriteMatrix.Texture
			tr.Origin = spriteMatrix.Origin
			tr.Dest = rl.Rectangle{
				Width:  currentAnimationFrame.Width,
				Height: currentAnimationFrame.Height,
			}
			tr.Frame = currentAnimationFrame
		}
		return true
	})
}
func (s *TextureRenderMatrixSystem) Destroy() {}
