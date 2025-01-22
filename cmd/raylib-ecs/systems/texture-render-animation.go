/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp_game/cmd/raylib-ecs/components"
	"gomp_game/pkgs/gomp/ecs"
)

// TextureRenderPosition is a system that sets Position of textureRender
type trAnimationController struct{}

func (s *trAnimationController) Init(world *ecs.World)        {}
func (s *trAnimationController) FixedUpdate(world *ecs.World) {}
func (s *trAnimationController) Update(world *ecs.World) {
	// Get component managers
	animations := components.AnimationService.GetManager(world)
	textureRenders := components.TextureRenderService.GetManager(world)

	// Update sprites and spriteRenders
	textureRenders.AllParallel(func(entity ecs.EntityID, tr *components.TextureRender) bool {
		if tr == nil {
			return true
		}

		animation := animations.Get(entity)
		if animation == nil {
			return true
		}

		frame := &tr.Frame
		if animation.Vertical {
			frame.Y += frame.Height * float32(animation.Current)
		} else {
			frame.X += frame.Width * float32(animation.Current)
		}

		return true
	})
}
func (s *trAnimationController) Destroy(world *ecs.World) {}
