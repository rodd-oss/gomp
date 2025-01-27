/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"fmt"
	"gomp_game/cmd/raylib-ecs/components"
	"gomp_game/pkgs/gomp/ecs"
	"time"

	"github.com/negrel/assert"
)

type animationPlayerController struct{}

func (s *animationPlayerController) Init(world *ecs.World) {}
func (s *animationPlayerController) Update(world *ecs.World) {
	animationPlayers := components.AnimationPlayerService.GetManager(world)

	animationPlayers.AllDataParallel(func(animation *components.AnimationPlayer) bool {
		animation.ElapsedTime += time.Duration(float32(world.DtUpdate().Microseconds())*animation.Speed) * time.Microsecond

		assert.True(animation.FrameDuration > 0, fmt.Errorf("frame duration must be greater than 0 (got %v)", animation.FrameDuration))

		// Check if animation is playing backwards
		if animation.Speed < 0 {
			for animation.ElapsedTime <= 0 {
				animation.ElapsedTime += animation.FrameDuration
				animation.Current--

				if animation.Current < animation.First {
					if animation.Loop {
						animation.Current = animation.Last
					} else {
						animation.Current = animation.First
					}
				}
			}
		} else {
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
		}

		return true
	})
}
func (s *animationPlayerController) FixedUpdate(world *ecs.World) {}
func (s *animationPlayerController) Destroy(world *ecs.World)     {}
