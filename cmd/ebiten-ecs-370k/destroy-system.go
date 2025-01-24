/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	"gomp_game/pkgs/gomp/ecs"
	"image/color"
)

type destroySystem struct {
	transformComponent *ecs.ComponentManager[transform]
	healthComponent    *ecs.ComponentManager[health]
	colorComponent     *ecs.ComponentManager[color.RGBA]
	movementComponent  *ecs.ComponentManager[movement]
	destroyComponent   *ecs.ComponentManager[empty]

	n int
}

func (s *destroySystem) Init(world *ecs.World) {
	s.transformComponent = transformComponentType.GetManager(world)
	s.healthComponent = healthComponentType.GetManager(world)
	s.colorComponent = colorComponentType.GetManager(world)
	s.movementComponent = movementComponentType.GetManager(world)
	s.destroyComponent = destroyComponentType.GetManager(world)

}
func (s *destroySystem) Run(world *ecs.World) {
	s.n = 0
	s.destroyComponent.All(func(e ecs.Entity, h *empty) bool {
		world.DestroyEntity(e)
		entityCount--

		return true
	})
}
func (s *destroySystem) Destroy(world *ecs.World) {}
