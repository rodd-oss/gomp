/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/examples/raylib-ecs/components"
	ecs2 "gomp/pkg/ecs"
)

// TextureRenderSprite is a system that prepares SpriteSheet to be rendered
type trSpriteMatrixController struct{}

func (s *trSpriteMatrixController) Init(world *ecs2.EntityManager)        {}
func (s *trSpriteMatrixController) FixedUpdate(world *ecs2.EntityManager) {}
func (s *trSpriteMatrixController) Update(world *ecs2.EntityManager) {
	// Get component managers
	spriteMatrixes := components.SpriteMatrixService.GetManager(world)
	textureRenders := components.TextureRenderService.GetManager(world)
	animationStates := components.AnimationStateService.GetManager(world)

	// Update sprites and spriteRenders
	spriteMatrixes.AllParallel(func(entity ecs2.Entity, spriteMatrix *components.SpriteMatrix) bool {
		if spriteMatrix == nil {
			return true
		}

		animationState := animationStates.Get(entity)
		if animationState == nil {
			return true
		}

		currentAnimationFrame := spriteMatrix.Animations[*animationState].Frame

		tr := textureRenders.Get(entity)
		if tr == nil {
			// Create new spriteRender
			newRender := components.TextureRender{
				Texture: spriteMatrix.Texture,
				Origin:  spriteMatrix.Origin,
				Frame:   currentAnimationFrame,
				Dest: rl.Rectangle{
					Width:  currentAnimationFrame.Width,
					Height: currentAnimationFrame.Height,
				},
			}

			textureRenders.Create(entity, newRender)
		} else {
			// Update spriteRender
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
func (s *trSpriteMatrixController) Destroy(world *ecs2.EntityManager) {}
