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

type hpSystem struct {
	transformComponent ecs.WorldComponents[transform]
	healthComponent    ecs.WorldComponents[health]
	colorComponent     ecs.WorldComponents[color.RGBA]
	movementComponent  ecs.WorldComponents[movement]
	destroyComponent   ecs.WorldComponents[empty]
}

func (s *hpSystem) Init(world *ecs.World) {
	s.transformComponent = transformComponentType.Instances(world)
	s.healthComponent = healthComponentType.Instances(world)
	s.colorComponent = colorComponentType.Instances(world)
	s.movementComponent = movementComponentType.Instances(world)
	s.destroyComponent = destroyComponentType.Instances(world)
}
func (s *hpSystem) Run(world *ecs.World) {
	s.healthComponent.AllParallel(func(entity ecs.EntityID, h *health) bool {
		h.hp--

		if h.hp <= 0 {
			s.destroyComponent.Set(entity, struct{}{})
		}

		return true
	})
}
func (s *hpSystem) Destroy(world *ecs.World) {}
