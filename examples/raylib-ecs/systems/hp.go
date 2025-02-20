/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package systems

import (
	"gomp/examples/raylib-ecs/components"
	ecs2 "gomp/pkg/ecs"
)

type hpController struct{}

func (s *hpController) Init(world *ecs2.EntityManager)   {}
func (s *hpController) Update(world *ecs2.EntityManager) {}
func (s *hpController) FixedUpdate(world *ecs2.EntityManager) {
	healths := components.HealthService.GetManager(world)

	healths.All(func(entity ecs2.Entity, h *components.Health) bool {
		h.Hp--

		if h.Hp <= 0 {
			world.Delete(entity)
		}

		return true
	})
}
func (s *hpController) Destroy(world *ecs2.EntityManager) {}
