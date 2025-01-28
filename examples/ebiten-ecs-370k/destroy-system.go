/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	ecs2 "gomp/pkg/ecs"
	"image/color"
)

type destroySystem struct {
	transformComponent *ecs2.ComponentManager[transform]
	healthComponent    *ecs2.ComponentManager[health]
	colorComponent     *ecs2.ComponentManager[color.RGBA]
	movementComponent  *ecs2.ComponentManager[movement]
	destroyComponent   *ecs2.ComponentManager[empty]

	n int
}

func (s *destroySystem) Init(world *ecs2.World) {
	s.transformComponent = transformComponentType.GetManager(world)
	s.healthComponent = healthComponentType.GetManager(world)
	s.colorComponent = colorComponentType.GetManager(world)
	s.movementComponent = movementComponentType.GetManager(world)
	s.destroyComponent = destroyComponentType.GetManager(world)

}
func (s *destroySystem) Run(world *ecs2.World) {
	s.n = 0
	s.destroyComponent.All(func(e ecs2.Entity, h *empty) bool {
		world.DestroyEntity(e)
		entityCount--

		return true
	})
}
func (s *destroySystem) Destroy(world *ecs2.World) {}
