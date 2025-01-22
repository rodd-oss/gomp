/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp_game/cmd/raylib-ecs/components"
	"gomp_game/pkgs/gomp/ecs"
	"time"
)

type animationController struct{}

func (s *animationController) Init(world *ecs.World) {}
func (s *animationController) Update(world *ecs.World) {
	animations := components.AnimationService.GetManager(world)

	animations.AllDataParallel(func(animation *components.Animation) bool {
		animation.ElapsedTime = animation.ElapsedTime + world.DtUpdate()*time.Duration(animation.Speed)

		for animation.ElapsedTime >= animation.FrameDuration {
			animation.ElapsedTime -= animation.FrameDuration
			animation.Current++

			if animation.Current > animation.Last {
				if animation.Loop {
					animation.Current = animation.First
				} else {
					animation.Current = animation.Last
				}
			}
		}

		return true
	})
}
func (s *animationController) FixedUpdate(world *ecs.World) {}
func (s *animationController) Destroy(world *ecs.World)     {}
