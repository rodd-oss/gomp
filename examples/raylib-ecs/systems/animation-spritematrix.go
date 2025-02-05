/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp/examples/raylib-ecs/components"
	"gomp/pkg/ecs"
	"time"
)

type animationSpriteMatrixController struct {
	selector ecs.Selector[struct {
		Player *components.AnimationPlayer
		State  *components.AnimationState
		Matrix *components.SpriteMatrix
	}]
}

func (s *animationSpriteMatrixController) Init(world *ecs.World) {
	world.RegisterSelector(&s.selector)
}
func (s *animationSpriteMatrixController) Update(world *ecs.World) {
	for c := range s.selector.All() {
		if c.Player.State == *c.State && c.Player.IsInitialized == true {
			continue
		}

		currentAnimation := c.Matrix.Animations[*c.State]

		c.Player.First = 0
		c.Player.Current = 0
		c.Player.Last = currentAnimation.NumOfFrames - 1
		c.Player.Loop = currentAnimation.Loop
		c.Player.Vertical = currentAnimation.Vertical
		c.Player.FrameDuration = time.Second / time.Duration(c.Matrix.FPS)
		c.Player.State = *c.State
		c.Player.Speed = 1
		c.Player.IsInitialized = true
	}
}
func (s *animationSpriteMatrixController) FixedUpdate(world *ecs.World) {}
func (s *animationSpriteMatrixController) Destroy(world *ecs.World)     {}
