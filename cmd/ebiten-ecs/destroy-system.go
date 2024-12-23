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
	transformComponent ecs.WorldComponents[transform]
	healthComponent    ecs.WorldComponents[health]
	colorComponent     ecs.WorldComponents[color.RGBA]
	movementComponent  ecs.WorldComponents[movement]
	destroyComponent   ecs.WorldComponents[empty]

	n int
}

func (s *destroySystem) Init(world *ecs.World) {
	s.transformComponent = transformComponentType.Instances(world)
	s.healthComponent = healthComponentType.Instances(world)
	s.colorComponent = colorComponentType.Instances(world)
	s.movementComponent = movementComponentType.Instances(world)
	s.destroyComponent = destroyComponentType.Instances(world)

}
func (s *destroySystem) Run(world *ecs.World) {
	s.n = 0
	s.destroyComponent.All(func(e ecs.EntityID, h *empty) bool {
		world.SoftDestroyEntity(e)
		entityCount--

		return true
	})

	world.Clean()
}
func (s *destroySystem) Destroy(world *ecs.World) {}
