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

type hpSystem struct {
	transformComponent *ecs2.ComponentManager[transform]
	healthComponent    *ecs2.ComponentManager[health]
	colorComponent     *ecs2.ComponentManager[color.RGBA]
	movementComponent  *ecs2.ComponentManager[movement]
	destroyComponent   *ecs2.ComponentManager[empty]
}

func (s *hpSystem) Init(world *ecs2.EntityManager) {
	s.transformComponent = transformComponentType.GetManager(world)
	s.healthComponent = healthComponentType.GetManager(world)
	s.colorComponent = colorComponentType.GetManager(world)
	s.movementComponent = movementComponentType.GetManager(world)
	s.destroyComponent = destroyComponentType.GetManager(world)
}
func (s *hpSystem) Run(world *ecs2.EntityManager) {
	s.healthComponent.AllParallel(func(entity ecs2.Entity, h *health) bool {
		h.hp--

		if h.hp <= 0 {
			s.destroyComponent.Create(entity, struct{}{})
		}

		return true
	})
}
func (s *hpSystem) Destroy(world *ecs2.EntityManager) {}
