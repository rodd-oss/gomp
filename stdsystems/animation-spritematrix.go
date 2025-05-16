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

func NewAnimationSpriteMatrixSystem() AnimationSpriteMatrixSystem {
	return AnimationSpriteMatrixSystem{}
}

type AnimationSpriteMatrixSystem struct {
	World            *ecs.EntityManager
	AnimationPlayers *stdcomponents.AnimationPlayerComponentManager
	AnimationStates  *stdcomponents.AnimationStateComponentManager
	SpriteMatrixes   *stdcomponents.SpriteMatrixComponentManager
}

func (s *AnimationSpriteMatrixSystem) Init() {}
func (s *AnimationSpriteMatrixSystem) Run() {
	s.AnimationPlayers.EachEntity(func(e ecs.Entity) bool {
		animationPlayer := s.AnimationPlayers.GetUnsafe(e)
		spriteMatrix := s.SpriteMatrixes.Get(e)
		if spriteMatrix == nil {
			return true
		}

		animationStatePtr := s.AnimationStates.GetUnsafe(e)
		if animationStatePtr == nil {
			return true
		}
		animationState := *animationStatePtr

		if animationPlayer.State == animationState && animationPlayer.IsInitialized == true {
			return true
		}

		currentAnimation := spriteMatrix.Animations[animationState]

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
