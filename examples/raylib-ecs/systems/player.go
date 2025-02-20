/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/examples/raylib-ecs/components"
	"gomp/examples/raylib-ecs/entities"
	"gomp/pkg/ecs"
)

type playerController struct {
	player entities.Player
}

func (s *playerController) Init(world *ecs.EntityManager) {
	s.player = entities.CreatePlayer(world)
	s.player.Position.X = 100
	s.player.Position.Y = 100

}
func (s *playerController) Update(world *ecs.EntityManager) {
	animationStates := components.AnimationStateService.GetManager(world)

	animationState := animationStates.Get(s.player.Entity)

	if rl.IsKeyDown(rl.KeySpace) {
		*animationState = entities.PlayerStateJump
	} else {
		*animationState = entities.PlayerStateIdle
		if rl.IsKeyDown(rl.KeyD) {
			*animationState = entities.PlayerStateWalk
			s.player.Position.X++
			s.player.Mirrored.X = false
		}
		if rl.IsKeyDown(rl.KeyA) {
			*animationState = entities.PlayerStateWalk
			s.player.Position.X--
			s.player.Mirrored.X = true
		}
		if rl.IsKeyDown(rl.KeyW) {
			*animationState = entities.PlayerStateWalk
			s.player.Position.Y--
		}
		if rl.IsKeyDown(rl.KeyS) {
			*animationState = entities.PlayerStateWalk
			s.player.Position.Y++
		}
	}
}
func (s *playerController) FixedUpdate(world *ecs.EntityManager) {}
func (s *playerController) Destroy(world *ecs.EntityManager)     {}
