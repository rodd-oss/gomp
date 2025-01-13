/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp_game/cmd/raylib-ecs/components"
	"gomp_game/pkgs/gomp/ecs"
)

type systemCalcHp struct{}

func (s *systemCalcHp) Init(world *ecs.World) {}
func (s *systemCalcHp) Run(world *ecs.World) {
	healths := components.HealthManager.Instances(world)

	healths.All(func(entity ecs.EntityID, h *components.Health) bool {
		h.Hp--

		if h.Hp <= 0 {
			world.DestroyEntity(entity)
		}

		return true
	})
}
func (s *systemCalcHp) Destroy(world *ecs.World) {}
