/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp_game/cmd/raylib-ecs/components"
	"gomp_game/cmd/raylib-ecs/entities"
	"gomp_game/pkgs/gomp/ecs"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type playerController struct {
	player entities.Player
}

func (s *playerController) Init(world *ecs.World) {
	s.player = entities.CreatePlayer(world)
	s.player.Position.X = 100
	s.player.Position.Y = 100
}
func (s *playerController) Update(world *ecs.World) {
	animations := components.AnimationService.GetManager(world)

	animation := animations.Get(s.player.Entity)

	if rl.IsKeyDown(rl.KeyD) {
		animation.Speed = 0
	} else {
		animation.Speed = 1
	}
}
func (s *playerController) FixedUpdate(world *ecs.World) {}
func (s *playerController) Destroy(world *ecs.World)     {}
