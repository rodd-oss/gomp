/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	rl "github.com/gen2brain/raylib-go/raylib"
	"gomp/examples/new-api/entities"
	"gomp/pkg/ecs"
	"gomp/stdcomponents"
	"time"
)

func NewPlayerSystem(
	world *ecs.World,
	spriteMatrixes *stdcomponents.SpriteMatrixComponentManager,
	positions *stdcomponents.PositionComponentManager,
	rotations *stdcomponents.RotationComponentManager,
	scales *stdcomponents.ScaleComponentManager,
	animationPlayers *stdcomponents.AnimationPlayerComponentManager,
	animationStates *stdcomponents.AnimationStateComponentManager,
	tints *stdcomponents.TintComponentManager,
	flips *stdcomponents.FlipComponentManager,
) *PlayerSystem {
	return &PlayerSystem{
		world:            world,
		spriteMatrixes:   spriteMatrixes,
		positions:        positions,
		rotations:        rotations,
		scales:           scales,
		animationPlayers: animationPlayers,
		animationStates:  animationStates,
		tints:            tints,
		flips:            flips,
	}
}

type PlayerSystem struct {
	world            *ecs.World
	player           entities.Player
	spriteMatrixes   *stdcomponents.SpriteMatrixComponentManager
	positions        *stdcomponents.PositionComponentManager
	rotations        *stdcomponents.RotationComponentManager
	scales           *stdcomponents.ScaleComponentManager
	animationPlayers *stdcomponents.AnimationPlayerComponentManager
	animationStates  *stdcomponents.AnimationStateComponentManager
	tints            *stdcomponents.TintComponentManager
	flips            *stdcomponents.FlipComponentManager
}

func (s *PlayerSystem) Init() {
	s.player = entities.CreatePlayer(s.world, s.spriteMatrixes, s.positions, s.rotations, s.scales, s.animationPlayers, s.animationStates, s.tints, s.flips)
	s.player.Position.X = 100
	s.player.Position.Y = 100
}
func (s *PlayerSystem) Run(dt time.Duration) {
	animationState := s.animationStates.Get(s.player.Entity)

	if rl.IsKeyDown(rl.KeySpace) {
		*animationState = entities.PlayerStateJump
	} else {
		*animationState = entities.PlayerStateIdle
		if rl.IsKeyDown(rl.KeyD) {
			*animationState = entities.PlayerStateWalk
			s.player.Position.X++
			s.player.Flip.X = false
		}
		if rl.IsKeyDown(rl.KeyA) {
			*animationState = entities.PlayerStateWalk
			s.player.Position.X--
			s.player.Flip.X = true
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
func (s *PlayerSystem) Destroy() {}
