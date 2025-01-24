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

type animatorController struct{}

func (s *animatorController) Init(world *ecs.World) {}
func (s *animatorController) Update(world *ecs.World) {
	animators := components.AnimatorService.GetManager(world)
	spritesheets := components.SpriteSheetService.GetManager(world)
	animations := components.AnimationService.GetManager(world)

	animators.AllParallel(func(ei ecs.Entity, a *components.Animator) bool {
		if a.LastState == a.State {
			return true
		}

		spritesheet := spritesheets.Get(ei)
		if spritesheet == nil {
			return true
		}

		animation := animations.Get(ei)
		if animation == nil {
			return true
		}

		a.LastState = a.State

		return true
	})
}
func (s *animatorController) FixedUpdate(world *ecs.World) {}
func (s *animatorController) Destroy(world *ecs.World)     {}
