/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/examples/new-api/components"
	"gomp/examples/new-api/entities"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
)

func NewPlayerSystem() PlayerSystem {
	return PlayerSystem{}
}

type PlayerSystem struct {
	EntityManager    *ecs.EntityManager
	Player           entities.Player
	SpriteMatrixes   *stdcomponents.SpriteMatrixComponentManager
	Positions        *stdcomponents.PositionComponentManager
	Rotations        *stdcomponents.RotationComponentManager
	Scales           *stdcomponents.ScaleComponentManager
	Velocities       *stdcomponents.VelocityComponentManager
	AnimationPlayers *stdcomponents.AnimationPlayerComponentManager
	AnimationStates  *stdcomponents.AnimationStateComponentManager
	Tints            *stdcomponents.TintComponentManager
	Flips            *stdcomponents.FlipComponentManager
	HP               *components.HealthComponentManager
	Controllers      *components.ControllerComponentManager
}

func (s *PlayerSystem) Init() {
	s.Player = entities.CreatePlayer(
		s.EntityManager, s.SpriteMatrixes, s.Positions, s.Rotations, s.Scales,
		s.Velocities, s.AnimationPlayers, s.AnimationStates, s.Tints, s.Flips,
	)

	s.Controllers.Create(s.Player.Entity, components.Controller{})

	s.Player.Position.X = 100
	s.Player.Position.Y = 100
}
func (s *PlayerSystem) Run() {
	animationState := s.AnimationStates.Get(s.Player.Entity)

	var speed float32 = 300

	s.Player.Velocity.X = 0
	s.Player.Velocity.Y = 0

	s.Controllers.AllData(func(c *components.Controller) bool {
		if rl.IsKeyDown(rl.KeySpace) {
			*animationState = entities.PlayerStateJump
		} else {
			*animationState = entities.PlayerStateIdle
			if rl.IsKeyDown(rl.KeyD) {
				*animationState = entities.PlayerStateWalk
				s.Player.Velocity.X = speed
				s.Player.Flip.X = false
			}
			if rl.IsKeyDown(rl.KeyA) {
				*animationState = entities.PlayerStateWalk
				s.Player.Velocity.X = -speed
				s.Player.Flip.X = true
			}
			if rl.IsKeyDown(rl.KeyW) {
				*animationState = entities.PlayerStateWalk
				s.Player.Velocity.Y = -speed
			}
			if rl.IsKeyDown(rl.KeyS) {
				*animationState = entities.PlayerStateWalk
				s.Player.Velocity.Y = speed
			}
		}

		if rl.IsKeyPressed(rl.KeyK) {
			s.EntityManager.Delete(s.Player.Entity)
		}
		return true
	})

}
func (s *PlayerSystem) Destroy() {}
