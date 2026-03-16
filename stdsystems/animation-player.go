/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package stdsystems

import (
	"gomp/pkg/ecs"
	"gomp/pkg/worker"
	"gomp/stdcomponents"
	"time"

	"github.com/negrel/assert"
)

func NewAnimationPlayerSystem() AnimationPlayerSystem {
	return AnimationPlayerSystem{}
}

type AnimationPlayerSystem struct {
	AnimationPlayers *ecs.ComponentManager[stdcomponents.AnimationPlayer]
	lastRunAt        time.Time
}

func (s *AnimationPlayerSystem) Init() {
	s.lastRunAt = time.Now()
}
func (s *AnimationPlayerSystem) Run() {
	dt := time.Since(s.lastRunAt)
	s.AnimationPlayers.ProcessComponents(func(animation *stdcomponents.AnimationPlayer, workerId worker.WorkerId) {
		animation.ElapsedTime += time.Duration(float32(dt.Microseconds())*animation.Speed) * time.Microsecond

		assert.True(animation.FrameDuration > 0, "frame duration must be greater than 0")

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

	})
	s.lastRunAt = time.Now()
}
func (s *AnimationPlayerSystem) Destroy() {}
