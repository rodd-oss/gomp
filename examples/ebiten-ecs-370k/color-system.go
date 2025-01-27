/*
This Source Code Form is subject to the terms of the Mozilla
Public License, v. 2.0. If a copy of the MPL was not distributed
with this file, You can obtain one at http://mozilla.org/MPL/2.0/.
*/

package main

import (
	ecs2 "gomp/pkgs/ecs"
	"image/color"
)

type colorSystem struct {
	transformComponent *ecs2.ComponentManager[transform]
	healthComponent    *ecs2.ComponentManager[health]
	colorComponent     *ecs2.ComponentManager[color.RGBA]
	movementComponent  *ecs2.ComponentManager[movement]

	baseColor color.RGBA
}

func (s *colorSystem) Init(world *ecs2.World) {
	s.transformComponent = transformComponentType.GetManager(world)
	s.healthComponent = healthComponentType.GetManager(world)
	s.colorComponent = colorComponentType.GetManager(world)
	s.movementComponent = movementComponentType.GetManager(world)

	s.baseColor = color.RGBA{25, 220, 200, 255}
}
func (s *colorSystem) Run(world *ecs2.World) {
	s.colorComponent.AllParallel(func(ei ecs2.Entity, c *color.RGBA) bool {
		health := s.healthComponent.Get(ei)
		if health == nil {
			return true
		}

		hpPercentage := float32(health.hp) / float32(health.maxHp)

		c.R = uint8(hpPercentage * float32(s.baseColor.R))
		c.G = uint8(hpPercentage * float32(s.baseColor.G))
		c.B = uint8(hpPercentage * float32(s.baseColor.B))
		c.A = s.baseColor.A

		return true
	})
}
func (s *colorSystem) Destroy(world *ecs2.World) {}
