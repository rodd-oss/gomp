/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package stdsystems

import (
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"time"
)

func NewAnimationSpriteMatrixSystem(world *ecs.World,
	animationPlayers *ecs.ComponentManager[stdcomponents.AnimationPlayer],
	animationStates *ecs.ComponentManager[stdcomponents.AnimationState],
	spriteMatrixes *ecs.ComponentManager[stdcomponents.SpriteMatrix]) *AnimationSpriteMatrixSystem {
	return &AnimationSpriteMatrixSystem{
		world:            world,
		animationPlayers: animationPlayers,
		animationStates:  animationStates,
		spriteMatrixes:   spriteMatrixes,
	}
}

type AnimationSpriteMatrixSystem struct {
	world            *ecs.World
	animationPlayers *ecs.ComponentManager[stdcomponents.AnimationPlayer]
	animationStates  *ecs.ComponentManager[stdcomponents.AnimationState]
	spriteMatrixes   *ecs.ComponentManager[stdcomponents.SpriteMatrix]
}

func (s *AnimationSpriteMatrixSystem) Init() {}
func (s *AnimationSpriteMatrixSystem) Run(dt time.Duration) {
	s.animationPlayers.AllParallel(func(e ecs.Entity, animationPlayer *stdcomponents.AnimationPlayer) bool {
		spriteMatrix := s.spriteMatrixes.Get(e)
		if spriteMatrix == nil {
			return true
		}

		animationStatePtr := s.animationStates.Get(e)
		if animationStatePtr == nil {
			return true
		}
		animationState := *animationStatePtr

		if animationPlayer.State == animationState && animationPlayer.IsInitialized == true {
			return true
		}

		currentAnimation := spriteMatrix.Animations[animationState]

		animationPlayer.First = 0
		animationPlayer.Current = 0
		animationPlayer.Last = currentAnimation.NumOfFrames - 1
		animationPlayer.Loop = currentAnimation.Loop
		animationPlayer.Vertical = currentAnimation.Vertical
		animationPlayer.FrameDuration = time.Second / time.Duration(spriteMatrix.FPS)
		animationPlayer.State = animationState
		animationPlayer.Speed = 1
		animationPlayer.IsInitialized = true

		return true
	})
}
func (s *AnimationSpriteMatrixSystem) Destroy() {}
